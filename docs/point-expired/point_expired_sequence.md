### Upsert Point ExpiredDate
```plantuml
actor User as user

participant appServer as app
database appDB as db

app -> app: env POINT_EXPIRED_DATE
user -> app: API pointAdd
app -> db: pointUpsert\nexpiredDate=today + POINT_EXPIRED_DATE  
```

### Exec Expired Point
```plantuml
participant Scheduler as sch
participant Worker as worker
participant AppServer as app

database AppDB as db
queue RabbitMQ as queue

sch -> db: find expired point users
db --> sch: user list

alt no expired users
    sch -> sch: finish
else expired users exist
    sch -> db: create job (PENDING)
    sch -> queue: publish userId
    sch -> sch: finish

    worker -> queue: consume message
    alt message received
        worker -> db: update job (RUNNING)
        worker -> app: expire point request (userId)
        app -> db: expire points (idempotent)
        app -> db: update job (DONE)
        app --> worker: success
        worker -> queue: ack
        worker -> worker: finish
    end
end
```