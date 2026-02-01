import logging
import time
from typing import Dict

from fastapi import APIRouter, FastAPI, HTTPException, Request, status
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import JSONResponse

from postfacta.core.environment import get_logging_level_from_env
from postfacta.routes.v0.handlers import router_v0


logger = logging.getLogger(__name__)
logging.basicConfig(
    format="[%(asctime)s][%(levelname)s] %(name)s: %(message)s",
    datefmt="%Y-%m-%d %H:%M:%S",
    level=get_logging_level_from_env(),
)


app = FastAPI(
    title="PostFacta",
    summary="Operational Incident Intelligence Service",
    contact={
        "name": "Chino Franco",
        "email": "chino.franco@gmail.com",
    },
)
startup_time = time.time()


@app.get("/", status_code=status.HTTP_200_OK, tags=["SYSTEM"])
def root():
    """Project main page."""
    return {"message": "Welcome to the PostFacta API!"}


@app.get("/healthz", status_code=status.HTTP_200_OK, tags=["SYSTEM"])
def health_check() -> Dict[str, str]:
    """Health check for the API."""
    return {"status": "healthy"}


@app.get("/service-info", status_code=status.HTTP_200_OK, tags=["SYSTEM"])
def service_info() -> Dict[str, object]:
    """Display the project information."""
    return {
        "uptime": time.time() - startup_time,
    }


@app.exception_handler(HTTPException)
async def http_exception_handler(request: Request, exc: HTTPException):
    """General exception handler."""
    return JSONResponse(
        status_code=exc.status_code,
        content={"status": exc.status_code, "message": exc.detail},
    )


api_base = APIRouter(prefix="/api")
api_base.include_router(router_v0)
app.include_router(api_base)
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],  # TODO: Adjust this to restrict origins
    allow_credentials=True,
    allow_methods=["GET"],
    allow_headers=["*"],
)
