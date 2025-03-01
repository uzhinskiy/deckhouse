---
title: "Модуль flant-integration"
---

Модуль выполняет функции по интеграции с различными сервисами Флант:
* Устанавливает в кластер madison-proxy в качестве alertmanager для Prometheus. Регистрируется в [Madison](#оповещения-в-madison).
* [Отправляет статистику](#статистика-о-состоянии-кластера), необходимую для расчета стоимости обслуживания кластера.
* [Отправляет логи](#логи-оператора-deckhouse) оператора Deckhouse, необходимые для облегчения процесса отладки.
* [Настраивает сбор метрик SLA](#метрики-sla).

## Сбор данных

### Куда Deckhouse отправляет данные?

Все данные отправляются через единую точку входа. Этой точкой является сервис `connect.deckhouse.io` (далее — *Connect* или *сервис Connect*).

Для авторизации при отправке данных каждый кластер Deckhouse отправляет ключ лицензии (license key) в качестве [Bearer Token](https://oauth.net/2/bearer-tokens/). Сервис Connect проверяет действительность ключа, и перенаправляет запрос в необходимый внутренний сервис Флант. 

Для отправки данных из кластера необходимо разрешить доступ до всех IP-адресов следующих DNS имен:
- `сonnect.deckhouse.io`;
- `madison-direct.deckhouse.io`.

### Какие данные отправляет Deckhouse?

> Узнайте о том, [как отключить отправку данных Deckhouse...](#как-отключить-отправку-данных-deckhouse)

Данные, отправляемые из кластера:
- Статистика о состоянии кластера: 
  - версия Kubernetes;
  - версия Deckhouse;
  - канал обновлений; 
  - количество узлов, и т.п.
- Оповещения, отправляемые в систему работы с инцидентами Madison.
- Метрики SLA по компонентам Deckhouse.
- Логи оператора Deckhouse.
- Способ подключения к master-узлам кластера.

#### Статистика о состоянии кластера

При помощи [shell-operator](https://github.com/flant/shell-operator) модуль flant-integration выполняет сбор метрик о состоянии объектов кластера. Затем с помощью [Grafana agent](https://github.com/grafana/agent) собранные метрики передаются по протоколу [Prometheus Remote Write](https://docs.sysdig.com/en/docs/installation/prometheus-remote-write/).

На основании собранных данных рассчитывается стоимость услуги [Managed Kubernetes](https://flant.ru/services/managed-kubernetes-as-a-service).

Среднее количество sample’ов отсылаемых одним кластером: **35 строк каждые 30 секунд**.

Пример собранных данных:
![](../../images/600-flant-integration/image1.png)
![](../../images/600-flant-integration/image2.png)

Дополнительно к метрикам собираемым модулем flant-integration, с помощью модуля [upmeter](../500-upmeter/) производится сбор [метрик доступности](#метрики-sla) для анализа выполнения SLA.

#### Оповещения в Madison

Madison — сервис обработки оповещений в составе платформы мониторинга компании Флант. Madison может принимать уведомления в формате Prometheus. 

При создании нового кластера Deckhouse:
1. При помощи ключа лицензии происходит автоматическая регистрация кластера в Madison.
2. Madison выдаёт кластеру ключ, необходимый для отправки уведомлений.
3. При помощи DNS-запроса для домена `madison-direct.flant.com` Deckhouse находит все доступные на текущий момент IP-адреса Madison.
4. Для каждого адреса создается Pod `madison-proxy`, в который Prometheus отправляет уведомления.

Схема отправки оповещений из кластера в Madison:
![](../../images/600-flant-integration/image3.png)

В среднем скорость отправки оповещений из кластера равняется **2 kb/s**. Но здесь следует учитывать, что чем больше инцидентов произошло в кластере, тем больше данных будет отправлено.

#### Логи оператора Deckhouse

Оператор Deckhouse является центральным компонентом всего кластера. Чтобы иметь данные необходимые для диагностики проблем в кластере, модуль flant-pricing настраивает модуль [log-shipper](../460-log-shipper/) для отправки логов в хранилище Loki компании Флант (не напрямую, а также через сервис Connect).

В логах содержится информация **только о компонентах Deckhouse, и нет секретных данных** касающихся кластера (примеры сообщений приведены на рисунке ниже). Собранная информация помогает понять, какие действия выполнял оператор Deckhouse, когда он их выполнял и с каким результатом.

Пример логов оператора Deckhouse:
![](../../images/600-flant-integration/image4.png)

При смене релиза Deckhouse, количество отправляемых данных логов достигает в среднем **150 записей в минуту**. При обычной работе, количество отправляемых данных логов достигает в среднем **20 записей в минуту**.

#### Метрики SLA

Модуль flant-pricing настраивает модуль [upmeter](../500-upmeter/) для отправки метрик, которые позволяют Флант контролировать выполнение условий соглашения об уровне сервиса (SLA) на компоненты кластера и компоненты Deckhouse.

### Как отключить отправку данных Deckhouse? 

Чтобы отключить регистрацию в Madison и отправку данных, необходимо отключить модуль `flantIntergation`.

**Важно!** Необходимо **обязательно** отключить модуль `flant-integration` в следующих случаях:
- В **тестовом кластере**, развернутом для экспериментов и т.п. Это правило не относится к кластерам разработки и тестовым кластерам, от которых нужно получать алерты.
- В любых **кластерах снятых с поддержки** Флант.

## Как осуществляется расчет стоимости?

Для каждой NodeGroup, за исключением выделенных мастеров, автоматически вычисляется тип биллинга. Существуют следующие
типы биллинга узлов:

* Ephemeral — если узел относится к NodeGroup с типом Cloud, то она автоматически относится к Ephemeral.
* VM — данный тип проставляется автоматически, если для узла удалось определить тип виртуализации с помощью команды
  [virt-what](https://people.redhat.com/~rjones/virt-what/).
* Hard — все остальные узлы автоматически относятся к данному типу.
* Special — данный тип необходимо вручную проставлять на NodeGroup, сюда относятся выделенные узлы, которые нельзя
  "потерять".

В случае, если в кластере есть узлы с типом биллинга Special или автоматическое определение сработало некорректно,
то вы всегда можете вручную установить корректный тип биллинга.

Для установки типа биллинга на узлах рекомендуется устанавливать аннотацию на NodeGroup, к которой относится узел:

```shell
kubectl patch ng worker --patch '{"spec":{"nodeTemplate":{"annotations":{"pricing.flant.com/nodeType":"Special"}}}}' --type=merge
```

Если в рамках одной NodeGroup есть узлы с разными типами биллинга, то можно навесить аннотацию отдельно на каждый объект Node:

```shell
kubectl annotate node test pricing.flant.com/nodeType=Special
```

### Определение статусов terraform-стейтов

Модуль опирается на метрики экспортируемые компонентом `terraform-exporter`. В них содержатся статусы соответствия
ресурсов в облаке/кластере с заданными в конфигурациях `*-cluster-configuration`.

#### Исходные метрики `terraform-exporter` и их статусы

- `candi_converge_cluster_status` соответствие конфигурации базовой инфраструктуры:
  - `error` — ошибка обработки, подробности смотреть в логе экспортера.
  - `destructively_changed` — `terraform plan` предполагает изменение объектов в облаке с удалением какого-либо из них.
  - `changed` — `terraform plan` предполагает изменение объектов в облаке без их удаления.
  - `ok`.
- `candi_converge_node_status` — соответствие конфигурации отдельных Node:
  - `error` — ошибка обработки, подробности смотреть в логе экспортера.
  - `destructively_changed` — `terraform plan` предполагает изменение объектов в облаке с удалением какого-либо из них.
  - `abandoned` — в кластере лишняя Node.
  - `absent` — в кластере не хватает Node.
  - `changed` — `terraform plan` предполагает изменение объектов в облаке без их удаления.
  - `ok`.
- `candi_converge_node_template_status` — соответствие `nodeTemplate` для `master` и `terranode` NodeGroup:
  - `absent` — NodeGroup отсутствует в кластере.
  - `changed` — параметры `nodeTemplate` расходятся.
  - `ok`.

#### Конечные метрики модуля `flant-integration` и механизм их получения

> Если модуль `terraform-manager` выключен в кластере — статус во всех метриках будет `none`. Данный статус следует трактовать как: стейта в кластере нет, но и не должно быть.

- Статус кластера (базовой инфраструктуры):
  - Используется значение метрики `candi_converge_cluster_status`.
  - В случае отсутствия метрики — `missing`.
- Статус `master` NodeGroup:
  - Берется "худший" статус из метрик `candi_converge_node_status` и `candi_converge_node_template_status` для `ng/master`.
  - В случае отсутствия обеих метрик — `missing`.
- Отдельный статус по каждой `terranode` NodeGroup:
  - Берется "худший" статус из метрик `candi_converge_node_status` и `candi_converge_node_template_status` для `ng/<nodeGroups[].name>`.
- Суммарный статус для всех `terranode` NodeGroup:
  - Берется "худший" статус из статусов, полученных для всех `terranode` NodeGroup.

> Статус `missing` так же будет фигурировать в конечных метриках, если `terraform-exporter` начнёт отдавать в своих метриках не описанные в модуле статусы. Иными словами статус `missing` это еще и некоего рода `fallback`-статус для ситуации, когда что-то пошло не так с определением "худшего" статуса.

#### Как определяется "худший" статус?

Мы считаем "худший" с точки зрения возможности автоматического применения существующих изменений.

Выбирается он по приоритету из следующей таблицы известных статусов:

| Статус                | Описание                                                                                  |
| --------------------- | ----------------------------------------------------------------------------------------- |
| error                 | Ошибка обработки стейта `terraform-exporter`'ом, подробности в его логе.                  |
| destructively_changed | `terraform plan` предполагает изменение объектов в облаке с удалением какого-либо из них. |
| abandoned             | В кластере лишняя Node.                                                                   |
| absent                | В кластере не хватает Node или NodeGroup.                                                 |
| changed               | `terraform plan` предполагает изменение объектов в облаке без их удаления.                |
| ok                    | Расхождений не обнаружено.                                                                |

