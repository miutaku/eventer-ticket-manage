# IMPORTANT

WIPです。
- Forwarder (cloudflare workerで実装でもいいかな、とも。)
- Parser
- Frontend (優先度低め)

がまだ出来上がっていません。
メールのサンプルを収集中のため、ForwarderとParserがWIPになっている感じ。


# eventer-ticket-manage

## budge

* draw.io

[![Draw.io to PNG](https://github.com/miutaku/eventer-ticket-manage/actions/workflows/drawio.yml/badge.svg)](https://github.com/miutaku/eventer-ticket-manage/actions/workflows/drawio.yml)

* dev

[![deploy cf-email-forwarder to dev](https://github.com/miutaku/eventer-ticket-manage/actions/workflows/dev.yml/badge.svg)](https://github.com/miutaku/eventer-ticket-manage/actions/workflows/deploy-cf-email-forwarder-dev.yml)
[![build and deploy to cloudrun as dev](https://github.com/miutaku/eventer-ticket-manage/actions/workflows/dev.yml/badge.svg)](https://github.com/miutaku/eventer-ticket-manage/actions/workflows/dev-cloudrun.yml)

# 概要

受け取ったチケットサイトのメールを元に、保持している・申し込んだチケットを一元管理するシステムです。
メールでのみではなく、手動でユーザーが登録することも可能としたい(WIP)。

# アーキテクチャ

![アーキテクチャ図](./infra-chart.png)
