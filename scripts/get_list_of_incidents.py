from typing import Final

from requests import Session


SERVER_URL: Final[str] = "http://localhost:8000"


def get_all_incidents_from_server() -> None:
    """Retrieve the list of incidents from the server."""
    with Session() as session:
        session.headers.update({"Content-Type": "application/json"})
        response = session.get(f"{SERVER_URL}/api/v0/incidents")
        assert response.status_code == 200
        incidents = response.json().get("incidents")
        print(f"Successfully retrieved {len(incidents)} incidents")
        for incident_id, incident_data in incidents.items():
            reporter = incident_data.get("reporter", "Unknown")
            print(f"Incident ID {incident_id}: Reported by {reporter}")


if __name__ == "__main__":
    get_all_incidents_from_server()
