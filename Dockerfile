FROM python:buster

WORKDIR /app

RUN pip install fastapi==0.68.0 uvicorn==0.14.0 aiofiles==0.7.0 aiojobs==0.3.0 async-timeout==3.0.1 RPi.GPIO

COPY rpi_temp_checker/ rpi_temp_checker/

CMD python -m rpi_temp_checker
