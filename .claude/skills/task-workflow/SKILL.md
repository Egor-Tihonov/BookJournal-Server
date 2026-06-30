---
name: task-workflow
description: Guided multi-phase workflow for non-trivial coding tasks aligned with the project CLAUDE.md operating model. Activate when user says "давай следующую таску", "возьмёмся за <X>", "продолжим таску <Y>", "что у меня по задачам", "/task-sel", or asks to start/resume any tracked task. Skill maintains tasks/<id>/{focus.md, result.md, SPEC.md, PLAN.md, DECISIONS.md, progress.md, RETROSPECTIVE.md} as source of truth, survives context compaction, runs 3-agent plan review (simplifier, architect, techlead), and delegates implementation to worker subagents. Do NOT activate for trivial tasks (single file, no architectural decisions, no new logic — see Lightweight Mode in CLAUDE.md), one-shot questions, quick debugging, or simple edits that do not need tracking.
---

# Task Workflow

## Цель

Провести пользователя через нетривиальную задачу с минимумом потерь контекста: от выбора таски до commit'а. Журнал в `tasks/<id>/` — источник правды, который выдерживает compaction и смену дня.

Скилл реализует **полный pipeline** из проектного `CLAUDE.md` (orchestrator + subagents + 7-файловая модель + retrospective). Не применяется к trivial задачам — для них действует Lightweight Mode из `CLAUDE.md`.

## Когда активируюсь

**Активация:**
- "давай следующую таску" / "возьмём новую задачу"
- "берёмся за MDT-1086" / "возьмёмся за <название>"
- "продолжим таску <id>" / "вернёмся к <id>"
- "что у меня в Asana" / "покажи мои таски"
- "/task-sel" (алиас если пользователь напишет)

**НЕ активируюсь:**
- Trivial задачи по классификации `CLAUDE.md` «Workflow calibration» + `AGENTS.md` `<rule id="workflow">`: один файл, без архитектурных решений, без новой логики, без protected zones / contracts / migrations. Опечатки, конфиги, документация, форматирование — всё это lightweight mode, без папки `tasks/<id>/`.
- Одноразовые вопросы ("как работает X?")
- Дебаг одного бага без трекинга
- Пользователь явно сказал "без таски / без журнала"

**При неуверенности:** спрашиваю "это полноценная таска с журналом или trivial-правка в lightweight mode?". Не эскалирую и не понижаю молча.

## Источник правды

Всегда:
1. Читаю `tasks/<id>/DECISIONS.md` первым при активации — там зафиксированы принятые решения и текущая фаза.
2. Текущая фаза — в шапке `DECISIONS.md`: строка `**Текущая фаза:** N — <имя>`.
3. При compaction / на следующий день — восстанавливаюсь полностью из файлов (восстановление — по `AGENTS.md` `<rule id="task-docs">`: источник правды это файлы).
4. Никогда не спрашиваю пользователя "что мы делали вчера" — это в файлах.

### Файлы таски (7 штук, по `AGENTS.md` `<rule id="task-docs">`)

| Файл | Роль |
|---|---|
| `focus.md` | Текущее направление внимания, scope, приоритеты, активные гипотезы, открытые вопросы. |
| `result.md` | Желаемое end-state: что должно быть в итоге; инварианты; неприемлемые исходы. |
| `SPEC.md` | Понятная спецификация изменения: что меняется, для кого, зачем, scope, acceptance criteria, TDD-матрица. |
| `PLAN.md` | Маршрут исполнения: упорядоченные bricks, назначения worker'ам, checkpoints, верификация. |
| `DECISIONS.md` | Принятые durable решения, отвергнутые альтернативы, текущая фаза. |
| `progress.md` | Текущее состояние исполнения: done / current / next / checks / gaps / blockers. |
| `RETROSPECTIVE.md` | Пост-задачный анализ ошибок и предложения по улучшению operating-model. |

Файлы subagent-вызовов — в `tasks/<id>/subagents/` по схеме `NNN-<role>-<model-tier>-<operation>-<result-type>-<scope>.md` (схема — Правило 11 ниже).

## Фазы (по `AGENTS.md` `<rule id="workflow">`)

### Phase 1 — Bootstrap
**Цель:** собрать стартовый контекст.

1. Если пользователь дал ID/имя — фиксирую. Иначе через Asana MCP вызываю `get_my_tasks(completed_since=now)`, показываю список (название, ID, дедлайн, проект), жду выбора.
2. Если у пользователя нет Asana-таски — прошу короткое английское имя (plain mode).
3. Через Asana MCP: `get_task(task_id, include_attachments=true)` — полное описание, parent, комментарии. `get_attachments(parent=task_id)` — список вложений.
4. Если во вложениях файлы — прошу пользователя приложить содержимое в чат или дать путь на диске.
5. Читаю `CLAUDE.md` и `AGENTS.md` проекта — это правила игры.
6. Читаю README, релевантные файлы (через Serena для Go, иначе Read/Grep), `.gitlab-ci.yml`, manifests.
7. Подтверждаю контекст пользователю: "вот что я понял из задачи. Верно?"

**Создание папки:**
- После подтверждения контекста — создаю `tasks/<id>/` (ID из Asana или короткое имя).
- 7 файлов из шаблонов (`.claude/skills/task-workflow/templates/`).
- Заполняю что уже известно; остальное — placeholder.
- **Проверяю** `.gitignore` проекта: если `tasks/` нет — добавляю с подтверждением пользователя.
- В `DECISIONS.md` фиксирую: **Текущая фаза: 1 — Bootstrap (завершено)**.

### Phase 2 — Direction
**Цель:** зафиксировать направление и destination до того как уйдём в реализацию.

1. Заполняю `focus.md`: attention, scope, out of scope, priorities, constraints, current hypotheses, open questions.
2. Заполняю `result.md`: desired end-state, expected behavior, invariants, unacceptable outcomes, preserved constraints, user-visible outcome, technical outcome.
3. Показываю оба файла пользователю: "вот куда идём, верно?"
4. Цикл уточнений пока пользователь не подтвердит.
5. Обновляю фазу в `DECISIONS.md`.

**Опционально**: если проблема framing'а нетривиальна — запускаю `plan-architect` или `plan-techlead` на проверку `focus.md` + `result.md` (analyzer checkpoint по `AGENTS.md` `<rule id="review-gates">`).

### Phase 3 — Specification
**Цель:** написать понятную спеку изменения.

1. Заполняю `SPEC.md` из шаблона: problem, intended change, scope, out of scope, current behavior, desired behavior, acceptance criteria, TDD/evidence matrix, risks and unknowns.
2. Для каждого acceptance criterion в TDD-матрице явно классифицирую: `automated-verifiable` или `manual-gap`.
3. Спека — это **что меняется и как успех узнать**, а не **как именно реализовать**. Архитектуру оставляю гибкой.
4. Показываю пользователю.

### Phase 4 — Planning
**Цель:** составить `PLAN.md` готовый к исполнению worker'ами.

1. Из шаблона `.claude/skills/task-workflow/templates/PLAN.md`.
2. Дроблю работу на bricks — каждый достаточно мал для одного worker-вызова.
3. Для каждого brick: scope, файлы, действие, ожидаемый output, TDD sequence (RED tests → impl), worker assignment, analyzer checkpoint.
4. Заполняю секции: Контекст, Шаги, Проверка, Риски, Что НЕ делаем, Rollback/compatibility.

**Anti-overengineering guard** — перед финализацией плана сам проверяю:
- Каждый новый тип — нужен?
- Каждое изменение existing сигнатуры — необходимо?
- Параметров не много?
- Тестов не чрезмерно?
- Коммиты не дроблятся без нужды?

Если сам нашёл overengineering — упрощаю до показа пользователю.

Обновляю фазу в `DECISIONS.md`.

### Phase 5 — Review (ОБЯЗАТЕЛЬНО)
**Цель:** прогнать план через analyzer subagents, минимизировать, проверить готовность.

**Принцип порядка:** сначала специалисты (architect, techlead) находят реальные риски и предлагают изменения. Затем **simplifier — ПОСЛЕДНИМ** — видит план уже после принятых правок и все их решения в `DECISIONS.md`, действует как финальный судья против overengineering. Это единственный проход simplifier'а на этой фазе.

Почему так: architect и techlead по роли расширяют (ищут риски → предлагают добавить тест, guard, абстракцию). Simplifier, запущенный **первым** — режет черновик без понимания рисков и может удалить нужное; запущенный **последним** — видит обоснованные расширения и решает, какие из них реально необходимы.

**Шаг 5.1.** Запускаю `plan-architect` и `plan-techlead` **параллельно** (одно сообщение с двумя Agent вызовами), модель — frontier. Каждый вызов сопровождаю файлом в `tasks/<id>/subagents/` по схеме `NNN-analyzer-frontier-plan-review-findings-<scope>.md`. Промпт примерно:
```
Прочитай tasks/<id>/{SPEC.md, PLAN.md, DECISIONS.md, focus.md, result.md}.
Найди <архитектурные проблемы | production risks> согласно твоей роли.
```

**Шаг 5.2.** Синтезирую находки обоих. Показываю пользователю компактно: принятое по умолчанию (no-brainer упрощения/исправления) + дискуссионное (где нужно решение пользователя).

**Шаг 5.3.** Обновляю `PLAN.md` с учётом принятых правок. В `DECISIONS.md` фиксирую решения раунда architect+techlead — это контекст для последующего simplifier'а.

**Шаг 5.4.** Запускаю `plan-simplifier` — **последним, с полным контекстом всего предыдущего**. Файл в `subagents/`: `NNN-analyzer-frontier-plan-review-findings-simplifier-<scope>.md`.

Промпт:
```
Это финальный раунд review перед реализацией. Прочитай:
- tasks/<id>/PLAN.md — текущая версия (после применения находок architect/techlead)
- tasks/<id>/DECISIONS.md — хронология решений, особенно последнее
  фиксирующее что добавил architect+techlead
- tasks/<id>/SPEC.md — acceptance criteria

Твоя роль — контр-баланс к architect/techlead, которые по природе расширяют план.
Найди в текущей версии плана overengineering, особенно среди того, что было добавлено
в последнем раунде:
- новые типы/абстракции без реального потребителя;
- тесты "на всякий случай" там, где семантика тривиальна;
- guard/fallback/normalization там, где simplicity > defensiveness;
- helper'ы-обёртки над одной строкой;
- scope creep ("попутно поправим X").

Доверяй чутью: каждое новое изменение должно доказать свою необходимость.
В конфликтах "архитектурная чистота vs меньше кода" — приоритет меньше кода,
если прод-риск отсутствует.

Формат findings — как обычно.
```

**Шаг 5.5.** Показываю simplifier-findings пользователю. Принятые — откатываю в `PLAN.md`. Приоритет simplifier'а сохраняется: если он говорит "откатить назад X" — предпочтение простоте, если не teardown реального прод-риска.

**Шаг 5.6.** КОНФЛИКТ-ПРАВИЛО:
- Если simplifier требует откатить что-то, что architect/techlead явно защищали как прод-риск — поднимаю конфликт:
  ```
  Simplifier: X можно убрать (обоснование).
  Techlead ранее: X закрывает риск Y (обоснование).
  Твоё решение?
  ```
- Пользователь решает. В случае "можно или так, или так" — simplicity.

**Шаг 5.7.** Финальная версия `PLAN.md`. Обновляю фазу в `DECISIONS.md`.

### Phase 6 — Execution
**Цель:** реализовать план через worker subagents под надзором primary agent.

**Default sequence для behavior changes (по `AGENTS.md` `<rule id="testing-evidence">`):**
1. Worker инспектирует текущее поведение и тестовые seams.
2. Worker пишет RED тесты (failing) на acceptance criteria из `SPEC.md`.
3. Primary agent проверяет: RED-evidence соответствует спеке.
4. Worker реализует минимальную production change.
5. Worker запускает targeted checks.
6. Analyzer (`plan-techlead` или подходящий) ревьюит diff и evidence.
7. Worker исправляет accepted findings.
8. Primary agent верифицирует финальное состояние.

**Делегирование worker'ам:**
- Каждый worker-вызов — отдельный файл в `tasks/<id>/subagents/` по схеме `NNN-worker-prev-implementation-patch-<scope>.md`.
- Модель worker'а — на тире ниже frontier (capable coding model, но дешевле).
- Brief содержит: assignment, context packet, expected output, TDD/evidence expectations, constraints, out of scope.
- Worker возвращает: changed files, commands and evidence, findings/blockers, handoff notes.

**Primary agent НЕ пишет product code** в Phase 6 — только координирует, ревьюит, читает diff'ы, запускает проверки, обновляет операционные документы. Исключение — trivial задача в Lightweight Mode (этот скилл туда не активируется, но если в процессе одна из bricks оказалась trivial, primary может выполнить её сам).

**После каждой brick:**
- Обновляю `progress.md`: Done / Current / Next / Checks / Gaps / Blockers.
- Если отклонение от плана — фиксирую в `progress.md` и в `DECISIONS.md` (что изменилось и почему).
- Даю **NEXT-блок**:
  ```
  ✓ Шаг N (<название>) готов.
  
  NEXT: Шаг N+1 — <название>
  Файл: <путь>, действие: <create/update>
  Worker: <модель/тир>
  
  Продолжаем?
  ```

### Phase 7 — Verification
**Цель:** проверить что работает целиком.

1. Запускаю команды из `PLAN.md` → Проверка: сборка, линтеры, тесты, race detector.
2. Прохожу поведенческий чеклист из `SPEC.md` / `PLAN.md`: либо автоматизированно (если покрыто тестами), либо прошу пользователя ручной прогон.
3. Если проблемы — возвращаюсь в Phase 6 на нужный шаг (с новой brick для worker'а).
4. `progress.md` → Checks — отмечаю checkboxes.
5. **Манчуальные gaps** — список из `SPEC.md` (`manual-gap`) явно фиксирую: что проверено руками, что осталось как accepted gap.

### Phase 8 — Retrospective и Final Report
**Цель:** извлечь уроки и закрыть таску.

**Retrospective loop** (по `AGENTS.md` `<rule id="task-docs">`, файл RETROSPECTIVE.md):
1. Запускаю если таска была meaningful (medium/large/multi-turn/risky/failed/heavily replanned). Пропускаю для тривиальных и проходных задач.
2. Заполняю `RETROSPECTIVE.md`: task outcome, what went wrong / right, agent mistakes, missed checkpoints, evidence gaps, reusable lessons, proposed operating-model changes.
3. Готовлю clean-context brief для analyzer'а под `subagents/NNN-analyzer-frontier-retrospective-review-findings-<scope>.md`. Контекст — финальный набор task-доков, релевантные subagent-файлы, summary diff, evidence, известные ошибки.
4. Запускаю analyzer (frontier model) на retrospective review. Возвращает findings по форме CLAUDE.md.
5. Применяю предложения только если они project-agnostic + предотвращают повторные failures + совместимы с моделью + не покрыты уже. Иначе — route в `CLAUDE.md`, `DECISIONS.md` или шаблоны.
6. Изменения в репозиторных `AGENTS.md`/`CLAUDE.md` (базовых документах operating-model) применяю **только** если текущая user-request это явно разрешает. Иначе — фиксирую как предложение и спрашиваю.

**Final Report:**
1. Финализирую `result.md`: что сделано, инварианты сохранены, user-visible outcome, technical outcome.
2. Финализирую `progress.md`: всё в Done.
3. Финализирую `DECISIONS.md`: **Текущая фаза: 8 — Retrospective + Final Report (завершено)**.
4. Summary пользователю — короткое резюме: что изменилось, какие subagents использовались, какие checks прошли, какие manual gaps остались, какие retrospective findings приняты/отложены, что осталось.
5. Спрашиваю: "коммитим? push'им?" — **не делаю без явного согласия**.

## Правила

### 1. NEXT-блок после каждой фазы
В конце моего сообщения на завершении фазы — явный блок:
```
✓ Фаза N (<название>) завершена.

NEXT: Фаза N+1 — <что>
<варианты / что предлагаю>

Поехали? / Твой вариант?
```
Пользователь не должен держать в голове "что дальше".

### 2. Никогда не создаю папку без подтверждения
Phase 1 создаёт `tasks/<id>/` только после "да, верно понял контекст".

### 3. Ручные прерывания
Если пользователь спросил что-то не по таске:
1. Отвечаю нормально на вопрос.
2. В конце:
   ```
   ↩ Вернёмся к таске <id>?
   Мы на Phase N (<название>). Остановились на <шаге/вопросе>.
   Следующий шаг: <X>
   Продолжаем или другое?
   ```

### 4. Смена таски
Если пользователь "бросаем текущую, берёмся за Z":
1. В `DECISIONS.md` текущей таски — секция "Статус": `Приостановлено <дата>. Текущая фаза: N — <имя>. Причина: <если сказали>.`
2. Обновляю `**Последнее обновление:** <дата>`.
3. Запускаю Phase 1 для Z.

### 5. После compaction
При первом сообщении в компактированной сессии по активной таске — следую `AGENTS.md` `<rule id="task-docs">` (источник правды — файлы):
1. Читаю `CLAUDE.md` и `AGENTS.md`.
2. Читаю `tasks/<id>/result.md` → SPEC.md → DECISIONS.md → PLAN.md → progress.md → focus.md → RETROSPECTIVE.md (если есть) → relevant subagents/*.
3. Инспектирую worktree.
4. Даю **сводку** пользователю:
   ```
   ↻ Восстановил контекст по tasks/<id>/:
   
   Задача: <кратко из SPEC.md>
   Фаза: N — <имя из DECISIONS.md>
   Остановились на: <progress.md → Current>
   Done: <progress.md → Done>
   Next: <progress.md → Next>
   
   Продолжаем с <X>?
   ```

### 6. Файлы = источник правды
- Не дублирую состояние в голове — всё в файлах.
- При любом сомнении — читаю файлы.
- Если файл устарел — обновляю его как первый шаг.
- Иерархия источника правды — по `AGENTS.md` `<rule id="precedence">`.

### 7. Шаблоны
Всегда из `.claude/skills/task-workflow/templates/`:
- `focus.md`, `result.md`, `SPEC.md`, `PLAN.md`, `DECISIONS.md`, `progress.md`, `RETROSPECTIVE.md`.

Никогда не копирую формат из старых `tasks/` — они могут быть устаревшими.

### 8. gitignore автопатч
При первой активации в новом проекте:
1. Проверяю наличие `.gitignore` в корне.
2. Если `tasks/` нет в игноре — спрашиваю: "Добавить `tasks/` в `.gitignore`? (это локальный журнал, не коммитится)".
3. С подтверждением — добавляю.

### 9. Работа без Asana
- Если у проекта нет Asana MCP или таски не в Asana — plain mode:
  - Прошу краткое английское имя для папки (`tasks/<short-name>/`).
  - Остальной workflow без изменений.

### 10. Ограничения реальности
- Не запускаю `git commit` / `git push` без явного согласия.
- Не применяю destructive изменения (reset --hard, force push, rm -rf) без явного разрешения.
- Если фаза требует действий пользователя (создать managed prompt в OpenAI, подставить credentials) — записываю в `result.md` чеклист для owner'а.

### 11. Subagent-файлы
Каждый вызов worker'а или analyzer'а сопровождается файлом в `tasks/<id>/subagents/`:
- Формат имени: `<NNN>-<role>-<model-tier>-<operation>-<result-type>-<scope>.md`
  - `role`: `worker` или `analyzer`
  - `model-tier`: `prev` (worker по умолчанию) или `frontier` (analyzer по умолчанию)
  - `operation`: `investigation`, `implementation`, `red-tests`, `plan-review`, `diff-review`, `result-review`, `retrospective-review`
  - `result-type`: `tests`, `patch`, `findings`, `evidence`, `decision`
  - `scope`: короткое kebab-case имя
- Структура файла самоописана в этом Правиле 11 (assignment / context / expected output / constraints → changed files / commands+evidence / findings / handoff).

## Антипаттерны (из опыта)

Сохранено как урок:

1. **Ревью-агенты усложняют план.** Три агента (architect + correctness + risk) по умолчанию находят проблемы в предложенном решении и чинят их — то есть добавляют код. Поэтому simplifier идёт **последним** и в одиночку, и его приоритет **выше** в конфликтах.

2. **Production-critical мышление для nice-to-have фич.** Не каждая фича достойна WaitGroup + retry + bounded pool + semaphore. Соразмеряю строгость гарантий с ценностью фичи.

3. **Вторжение в existing code.** Правило: новый код добавляется, existing не меняется. Если нужно менять сигнатуры/schema — проверяю есть ли injection-точка где этого избежать.

4. **Преждевременная обобщённость.** Struct без методов = часто просто набор параметров. Не ввожу типы ради "читаемости" — убеждаюсь что они упрощают.

5. **Избыточная observability.** 2-3 метрики достаточно для MVP. Action history, отдельные notes, колонки в UI — добавляю когда увидим реальную операционную потребность.

6. **Premature commit splitting.** Разбивка фичи на C1/C2a/C2b — когда это один feature flag и один rollout. Разбиваю только если (а) независимый deploy, (б) review становится проще.

7. **Trivial задача в полном pipeline.** Опечатка / правка одного env var не требует focus/result/SPEC/PLAN/DECISIONS/progress/RETROSPECTIVE + worker'ов. Это lightweight mode из CLAUDE.md — реализую сам, без папки `tasks/<id>/`. Если сомневаюсь между lightweight и full pipeline — спрашиваю.

## Чеклист перед завершением workflow

- [ ] `tasks/<id>/DECISIONS.md` финализирован (Текущая фаза: 8)
- [ ] `tasks/<id>/SPEC.md` отражает финальный scope
- [ ] `tasks/<id>/PLAN.md` финален (с учётом отклонений)
- [ ] `tasks/<id>/result.md` заполнен (end-state, инварианты, outcomes)
- [ ] `tasks/<id>/progress.md` — всё в Done, манчуальные gaps отмечены
- [ ] `tasks/<id>/RETROSPECTIVE.md` заполнен (если таска была meaningful)
- [ ] Все subagent-файлы в `tasks/<id>/subagents/` сохранены
- [ ] Все checkboxes в "Проверка" отмечены
- [ ] Если есть ops-шаги для owner'а (managed prompts, deploy, включение флагов) — в `result.md` отдельный чеклист "Шаги для owner"
- [ ] Пользователь подтвердил завершение
