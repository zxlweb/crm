# 会议 / 阶段笔记（Meeting Notes）

本目录存放**阶段小结、验收记录、阻塞与联调笔记**，不是任务勾选（见 [tasks/](../tasks/)），也不是 QA / Review 正式报告。

## 命名规范

| 场景 | 文件名 |
|------|--------|
| 按 Phase | `phase-N-notes.md` |
| 按会议 | `YYYY-MM-DD-{主题}.md` |
| 按模块迭代 | `{模块}-notes.md` |

## 文件索引

| 文件 | 说明 |
|------|------|
| [phase-0-notes.md](./phase-0-notes.md) | Phase 0 基础架构收尾笔记 |

## 协作规则

1. **任务勾选**只在 [tasks/00-mvp-task-breakdown.md](../tasks/00-mvp-task-breakdown.md)  
2. **QA 报告** → [docs/qa/](../qa/)  
3. **Code Review** → [docs/reviews/](../reviews/)  
4. 本目录写：验收要点、风险、本地命令、相关链接  
5. 更新后执行 `make docs-html`
