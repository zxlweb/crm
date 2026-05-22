# QA 测试报告

本目录存放**测试执行结果**（测了什么、通没通），不是任务进度。

## 命名规范

| 场景 | 文件名 | 模板 |
|------|--------|------|
| 按 Phase 收尾 | `phase-N-qa.md` | [qa-test-template.md](../templates/qa-test-template.md) |
| 按业务模块 | `{模块}-qa.md` | 同上 |

## 文件索引

| 文件 | 范围 | 结论 |
|------|------|------|
| [phase-0-qa.md](./phase-0-qa.md) | Phase 0 基础架构 | 通过 |

## 协作规则

1. 从 `docs/templates/qa-test-template.md` 复制新建，勿改模板正文  
2. 任务清单只在 [tasks/00-mvp-task-breakdown.md](../tasks/00-mvp-task-breakdown.md) 勾选  
3. 在 [meeting-notes/phase-N-notes.md](../meeting-notes/phase-N-notes.md) 中链接本目录报告  
4. 更新后执行 `make docs-html`
