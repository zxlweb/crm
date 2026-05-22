# Code Review 记录

本目录存放 **Code Review 正式结论**（过没过、改什么），不是任务进度。

## 命名规范

| 场景 | 文件名 | 模板 |
|------|--------|------|
| 按 Phase 收尾 | `phase-N-review.md` | [code-review-template.md](../templates/code-review-template.md) |
| 按业务模块 / PR | `{模块}-review.md` | 同上 |

## 文件索引

| 文件 | 范围 | 结论 |
|------|------|------|
| [phase-0-review.md](./phase-0-review.md) | Phase 0 基础架构 | 通过 |

## 协作规则

1. 从 `docs/templates/code-review-template.md` 复制新建  
2. Review 完成后在 [tasks/](../tasks/) 主清单勾选，并在 [meeting-notes/](../meeting-notes/) 更新链接  
3. 更新后执行 `make docs-html`
