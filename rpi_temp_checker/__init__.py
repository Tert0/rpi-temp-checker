"""RPi Temperature Checker"""
__version__ = "1.0.0"
from fastapi import FastAPI
import aiofiles
from enum import Enum
from os import getenv
import asyncio
import aiojobs
import RPi.GPIO as GPIO

app = FastAPI(title="Raspberry Pi Temperature Checker")


class Scheduler:
    scheduler = None

    async def init(self):
        self.scheduler = await aiojobs.create_scheduler()

    async def __call__(self):
        return self.scheduler


HIGHT_TEMP: float = float(getenv("HIGHT_TEMP", 55.0))
CRIT_TEMP: float = float(getenv("CRIT_TEMP", 70.0))
TIMEOUT: float = float(getenv("TIMEOUT", 10))
COOLER_PIN: int = int(getenv("COOLER_PIN", 3))


async def get_temp() -> float:
    async with aiofiles.open("/sys/class/hwmon/hwmon0/device/temp", "r") as file:
        content = await file.read()
    return int(content) / 1000


class TempStatus(str, Enum):
    NORMAL = "normal"
    HIGH = "high"
    CRITICAL = "critical"


async def get_status() -> TempStatus:
    temp = await get_temp()
    if temp >= CRIT_TEMP:
        return TempStatus.CRITICAL
    elif temp >= HIGHT_TEMP:
        return TempStatus.HIGH
    return TempStatus.NORMAL


@app.get("/temp")
async def get_temp_route():
    return await get_temp()


@app.get("/status", response_model=TempStatus)
async def get_status_route():
    return await get_status()


async def check_job(scheduler):
    status: TempStatus = await get_status()
    if status == TempStatus.NORMAL:
        GPIO.output(COOLER_PIN, GPIO.LOW)
    elif status == TempStatus.HIGH:
        GPIO.output(COOLER_PIN, GPIO.HIGH)
    elif status == TempStatus.CRITICAL:
        GPIO.output(COOLER_PIN, GPIO.HIGH)
        print("CRITICAL TEMP")
    await asyncio.sleep(TIMEOUT)
    await scheduler.spawn(check_job(scheduler))


@app.on_event("startup")
async def on_startup():
    GPIO.setmode(GPIO.BOARD)
    GPIO.setwarnings(False)
    GPIO.setup(COOLER_PIN, GPIO.OUT)
    scheduler = await aiojobs.create_scheduler()
    await scheduler.spawn(check_job(scheduler))

if __name__ == '__main__':
    import uvicorn
    uvicorn.run(app, host="8.8.8.8", port=8080)
