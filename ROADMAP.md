# Roadmap

## v0.1.x

目标：补强稳定性与接入体验。

- 扩大 live integration 覆盖面。
- 补更多核心 API 的 Go doc 和 examples。
- 收紧 token 暴露面与文档化安全边界。
- 增加更高层 façade API，减少调用方样板代码。

## v0.2.0

目标：增强类型安全与批量管理能力。

- 将更多协议字段从 magic string/int 提升为类型或枚举。
- 增加更完整的批量文件/目录管理能力。
- 增加上传进度、取消与重试能力。
- 收敛更结构化的错误模型。

## v1.0.0

目标：稳定公共 API，形成长期维护承诺。

- 明确兼容性与版本演进策略。
- 完整 examples / 文档体系。
- 更成熟的 CI / release / integration 回归矩阵。
- 发布稳定 API 面并冻结核心契约。
