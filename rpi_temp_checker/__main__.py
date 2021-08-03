from . import app
import uvicorn
uvicorn.run(app, host="8.8.8.8", port=8080)
