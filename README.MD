# kafka-stub
## Сервис-заглушка для кафки
[dockerHub](https://hub.docker.com/repository/docker/asavt92/kafka-stub)



Сервис вычитывает json сообщение из топика, валидирует его (json схема),  на основе маппинга полей пробрасывает необходимые поля в ответ, и отправляет в указанный топик.

### Конфигурирование
[Файл содержит дефолтные конфиги + описание](configs/config.yml)

Конфиги могут быть переопределены переменными окружения (UPPER_CASE).

### Как использовать
#### docker image

Образ содержит только файл с дефолтными конфигами. Ответные сообщения сервиса необходимо примонтировать по пути `out.json.examples.dir.path`.

json схему для валидации следует монтировать попути `in.json.schema.path` 

Запустить локально:

    docker run --rm --name=kafka-stub --network=host \
     -v ${PWD}/data:/data:ro \
     -v ${PWD}/configs:/configs:ro \
     -e KAFKA_GROUP=some_group \
        kafka-stub:latest
        
        
### TODO
- [ ] Добавить проброс значений полей из запроса в ответ
- [ ] Тесты
- [ ] Openshift manifest
- [ ] Makefile