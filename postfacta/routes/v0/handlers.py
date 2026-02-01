import logging

from fastapi import APIRouter, status

from postfacta.core.models import NewIncidentRequest, create_new_incident


logger = logging.getLogger(__name__)
router_v0 = APIRouter(prefix="/v0", tags=["V0"])


incidents_router = APIRouter(prefix="/incidents", tags=["INCIDENTS"])

@incidents_router.post("/start", status_code=status.HTTP_201_CREATED)
async def create_incident(incident_request: NewIncidentRequest):
    """Create a new incident."""
    new_incident = create_new_incident(incident_request)
    logger.info(f"New incident reported by {new_incident.reporter}: {new_incident.id}")
    return {"message": "Incident created", "id": new_incident.id}


router_v0.include_router(incidents_router)
