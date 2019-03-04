# credential-store
a handwriting credential store backed by postgres 

## developer guide
To run the program. do 
```
docker-compose up --build
```

this will bring up rabbitmq, postgres and credential-store all in dockers

To publish some message in the queue to trying the auth events, one can run
```
pip install -r ./scripts/python_publisher/requirements.txt
python ./scripts/python_publisher/auth_req_publish.py
```
This will send a create action, one success authentication action and one failed auth actions