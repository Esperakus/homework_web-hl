# Домашняя работа "Nginx - балансировка и отказоустойчивость"

По условиям задачи необходимо создать стенд, который будет обеспечивать:
- балансировку и отказоустойчивость бэкенда веб-приложения с помощью Nginx
- отказоустойчивость самих балансировщиков с помощью keepalived

Однако, при выполнении данной работы в среде Yandex Cloud я, да и наверное все остальные студенты данного курса, тут же столкнулись с проблемой: виртуальные сети Ya.Cloud не поддерживают VRRP, multicast и в целом L2. Соответственно отказоустойчивость с помощью "переезда" ip адреса в этих сетях обеспечить невозможно. :(

После нескольких вопросов в Slack в канале группы о том, как поступать, Алексей Цыкунов разрешил использовать в данной работе штатный Network Load Balancer от Yandex.Cloud.

Данный репозиторий содержит:

- Манифесты terraform для создания инфраструктуры проекта:
  - штатный балансировщик yandex.cloud, который будет проводить периодически health-check воркеров nginx и балансировать входящий трафик между ними
  - 2 воркера nginx, которые в свою очередь настроены на простейшую балансировку трафика на бэкенды веб приложения
  - 2 воркера бэкенда, на которых в systemd запущено простейшее приложение на go, слушающее порт 8090. При запросе отдаёт имя бэкенда (чтобы понять, на какой бэкенд прилетел запрос из Nginx) и версию БД
  - 1 экземпляр БД Postgresql 13 c базой test и пользователем БД test, чтоб принимать запросы от бэкенда
  - и, наконец, виртуалка с установленным ансиблем, чтоб развернуть вышеупомянутые роли. Выступает так же в роли Jump host проекта, поскольку единственная имеет внешний ip (не считая балансировщика yandex.cloud, он тоже имеет внешний ip)

- Роли ansible для приведения виртуальных машин в проекте в требуемое состояние.

При разворачивании стенда создаются ВМ с параметрами:
- 2 CPU;
- 2 GB RAM;
- 10 GB диск;
- операционная система CentOS 8 Stream;

Для разворачивания стенда необходимо:
1. Инициализировать рабочую среду Terraform:

```
$ terraform init
```
В результате будет установлен провайдер для подключения к облаку Яндекс.

2. Запустить разворачивание стенда:
```
$ terraform apply
```
В процессе разворачивания будут запрошены cloud_id, folder_id и iam-token. При желании эти значения можно задать соответсвующим переменным в variables.tf. В выходных данных будут показаны все внешние и внутренни ip адреса. Для проверки работы стенда необходимо в браузере или с помощью curl зайти на ip адрес балансировщика yandex.cloud, который можно посмотреть в выходных данных, например:

```
вывод  terraform apply:

...
external_ip_address_lb = tolist([
  {
    "external_address_spec" = toset([
      {
        "address" = "51.250.84.78"
        "ip_version" = "ipv4"
...
```
заходим на ip балансировщика (можно в сочетании с watch, чтоб наглядно видеть смену бэкендов):
```
curl http://51.250.84.78
```
в выводе увидим:
```
backend0.ru-central1.internal
PostgreSQL 13.9 on x86_64-pc-linux-gnu, compiled by gcc (GCC) 8.5.0 20210514 (Red Hat 8.5.0-15), 64-bit
```
или
```
backend1.ru-central1.internal
PostgreSQL 13.9 on x86_64-pc-linux-gnu, compiled by gcc (GCC) 8.5.0 20210514 (Red Hat 8.5.0-15), 64-bit
```
что говорит о том, что запрос может перенаправляться Nginx на разные бэкенды.

Можно зайти по ssh на джамп-хост, с которого можно попасть на любую ВМ внутри стенда. Для этого из рабочей папки проекта надо выполнить:

```
[homework02]$ ssh cloud-user@{external_ip_address_ansible} -i id_rsa
```
**external_ip_address_ansible** можно посмотреть в выводе terraform или консоли yandex.cloud.

С джамп-хоста можно ходить по ssh по всем машинам внутри проекта по их внутренним ip адресам или hostname (nginx0, nginx1, backend0, backend1, db).

Если останавливать по одной службы nginx на балансировщиках nginx0, nginx1 можно увидеть, что запросы всё равно идут к бэкендам благодаря работе балансировщика yandex.cloud.

Если останавливать по одной службы go_web.service на бэкендах backend0, backend1, можно увидеть, что запросы будут идти только на работающий бэкенд благодаря работе балансировки в nginx. Если запустить снова - оба бэкенда будут в работе спустя некоторое время.
