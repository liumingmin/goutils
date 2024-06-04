![gopher](gopher.png "goutils")


[![Go](https://github.com/liumingmin/goutils/actions/workflows/go.yml/badge.svg)](https://github.com/liumingmin/goutils/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/liumingmin/goutils)](https://goreportcard.com/report/github.com/liumingmin/goutils)
[![codecov](https://codecov.io/gh/liumingmin/goutils/graph/badge.svg?token=BQRDOY3CDX)](https://codecov.io/gh/liumingmin/goutils)
![GitHub last commit](https://img.shields.io/github/last-commit/liumingmin/goutils)
![GitHub Tag](https://img.shields.io/github/v/tag/liumingmin/goutils)
![GitHub License](https://img.shields.io/github/license/liumingmin/goutils)

**其他语言版本: [English](README.md), [中文](README_zh.md).**

# 简介

该仓库旨在为 Golang 开发人员提供一系列实用且易于使用的工具，帮助他们提高开发效率和工作效率。这些工具涵盖了各种领域，包括算法库、容器库、缓存工具、文件处理、Http与Websocket网络、NoSql数据库访问等。

# 仓库特色
- 方案定位: 立足于解决微服务框架之外繁琐工作，与开发框架形成差异化互补。
- 易于集成: 低耦合，便于与各种项目集成，如go-zero. 
- 功能聚焦: 工具库涵盖了各种常见功能的工具实现，为了避免重复造轮子没有开发自己的SQL ORM，更加灵活的选用开源社区强大的ORM方案，站在巨人的肩膀才能走的更远。


# 目录
- [算法模块](algorithm/README_zh.md)
- [缓存模块](cache/README_zh.md)
- [yaml配置模块](conf/README_zh.md)
- [容器模块](container/README_zh.md)
- [数据库](db/README_zh.md)
- [日志库](log/README_zh.md)
- [网络库](net/README_zh.md)
- [通用工具库](utils/README_zh.md)
- [websocket客户端和服务端库](ws/README_zh.md)

# 路线图

## 目标：

* 构建高质量、高性能的 Golang 项目
* 提高项目知名度和影响力

## 阶段划分：

### 第一阶段：基础建设

* **目标:** 确保代码质量和基础完善
* **关键指标:**
    - [x] 所有代码经过 `go test` 和 `codecover` 测试
    - [x] GitHub Star 超过 200
* **任务:**
    - [x] 编写单元测试和集成测试，确保代码功能正确性
    - [x] 使用 `codecover` 工具测量代码覆盖率，并持续提高覆盖率
    - [x] 撰写清晰的项目文档和 README 文件

### 第二阶段：性能优化和扩展

* **目标:** 提升项目性能和扩展性
* **关键指标:**
    * 主要代码经过性能测试，性能指标达到预期
    * 所有代码的代码覆盖率超过 80%
    * GitHub Star 超过 500
* **任务:**
    * 识别性能瓶颈，并进行针对性的优化
    * 使用负载测试工具评估项目性能，并持续改进
    * 完善代码结构和设计，提高代码可维护性和可扩展性
    * 发布新版本，并记录版本变更日志

### 第三阶段：社区运营和推广

* **目标:** 壮大项目社区，扩大项目影响力
* **关键指标:**
    * GitHub Star 超过 1000
    * 开设独立的官方站点
    * 积极参与开源社区活动
* **任务:**
    * 建立社区交流平台，例如论坛、QQ 群等
    * 组织线上线下技术交流活动，分享项目经验
    * 撰写博客文章、技术教程等，传播项目知识
    * 积极参与相关开源会议和活动，推广项目
