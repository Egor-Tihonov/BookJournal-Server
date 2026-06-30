# task-workflow skill

Guided 8-phase workflow for non-trivial coding tasks, aligned with the project
`CLAUDE.md` operating model (orchestrator + worker/analyzer subagents +
7-file task folder + retrospective).

## Phases (по `AGENTS.md` `<rule id="workflow">`)

1. **Bootstrap** — pick task (Asana or plain), gather attachments/docs, read `CLAUDE.md` + `AGENTS.md`, inspect repo, create `tasks/<id>/`
2. **Direction** — write `focus.md` + `result.md` (attention, scope, destination, invariants)
3. **Specification** — write `SPEC.md` (problem, intended change, acceptance criteria, TDD matrix)
4. **Planning** — write `PLAN.md` (bricks for workers, analyzer checkpoints, verification commands)
5. **Review** — 3 agents: `plan-architect` + `plan-techlead` in parallel, then `plan-simplifier` last with full context
6. **Execution** — delegate bricks to worker subagents (TDD: RED tests → implementation → checks), update `progress.md` after each
7. **Verification** — build, tests, race detector, behavioral checklist, manual gaps
8. **Retrospective + Final Report** — fill `RETROSPECTIVE.md` for meaningful tasks, finalize `result.md` and `progress.md`, ask about commit

## When NOT to activate

For trivial tasks (single file, no architectural decisions, no new logic — see
the `Workflow calibration` section in project `CLAUDE.md`), skip this skill and
implement directly. Опечатки, конфиги, форматирование, документация — это
lightweight mode, без папки `tasks/<id>/`.

## Files

- `SKILL.md` — full workflow specification + activation triggers + rules
- `templates/focus.md` — attention vector, scope, hypotheses
- `templates/result.md` — desired end-state, invariants, outcomes
- `templates/SPEC.md` — change specification, acceptance criteria, TDD matrix
- `templates/PLAN.md` — execution bricks, worker assignments, analyzer checkpoints
- `templates/DECISIONS.md` — durable decisions, supersessions, current phase
- `templates/progress.md` — done/current/next/checks/gaps/blockers
- `templates/RETROSPECTIVE.md` — post-task analysis and operating-model proposals

## Related agents

Three review agents in `.claude/agents/`:
- `plan-architect` — checks architectural fit with codebase (runs in parallel with techlead)
- `plan-techlead` — production readiness, regression risks, observability (runs in parallel with architect)
- `plan-simplifier` — finds overengineering, reduces scope (runs **last**, with full context of architect/techlead findings)

Order rationale: architect + techlead naturally expand the plan (find risks →
propose guards/tests/abstractions). Simplifier last sees the expanded plan and
cuts what is not justified.

## Source of truth

Files in `tasks/<id>/` survive context compaction. On activation, skill reads
`DECISIONS.md` (header has `Текущая фаза: N`) and `progress.md` to restore
state without asking the user "what were we doing".

Subagent invocations are persisted under `tasks/<id>/subagents/` using the
naming scheme `<NNN>-<role>-<model-tier>-<operation>-<result-type>-<scope>.md`
(see SKILL.md rule 11).

## Lessons baked in

From real experience (see antipattern section in SKILL.md):
- Review agents tend to **add** complexity, not reduce it → simplifier runs **last** with full context, simplicity wins ties
- Treat nice-to-have features with proportional rigor (not production-critical)
- Injection, not intrusion: new code added, existing code not changed
- 2-3 metrics enough for MVP observability, not 10
- One commit for one feature flag + one rollout
- Trivial tasks belong in lightweight mode (CLAUDE.md), not in the full pipeline
