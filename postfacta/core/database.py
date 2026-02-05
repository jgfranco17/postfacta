import logging
from abc import ABC, abstractmethod

from postfacta.core.errors import IncidentNotFoundError
from postfacta.core.models import Incident

logger = logging.getLogger(__name__)


class DataClient(ABC):
    """Defines the interface for database clients."""

    @abstractmethod
    def connect(self) -> None:
        """Establish a connection to the database."""
        raise NotImplementedError("Client needs implementation of this method")

    @abstractmethod
    def disconnect(self) -> None:
        """Close the connection to the database."""
        raise NotImplementedError("Client needs implementation of this method")

    @abstractmethod
    def register(self, incident: Incident) -> None:
        """Register a new incident in the database system."""
        raise NotImplementedError("Client needs implementation of this method")

    @abstractmethod
    def update_entry(self, incident: Incident) -> None:
        """Update an existing incident entry in the database."""
        raise NotImplementedError("Client needs implementation of this method")

    @abstractmethod
    def get_by_id(self, incident_id: str) -> Incident:
        """Retrieve an incident by its ID."""
        raise NotImplementedError("Client needs implementation of this method")

    @abstractmethod
    def get_all(self) -> dict[str, Incident]:
        """Retrieve all incidents from the database."""
        raise NotImplementedError("Client needs implementation of this method")

    @abstractmethod
    def remove_by_id(self, incident_id: str) -> None:
        """Remove an incident from the database by its ID."""
        raise NotImplementedError("Client needs implementation of this method")


class InMemoryClient(DataClient):
    """In-memory database client implementation."""

    def __init__(self) -> None:
        self._storage: dict[str, Incident] = {}

    def connect(self) -> None:
        """Simulate connecting to an in-memory database."""
        logging.info("Connected to in-memory database")

    def disconnect(self) -> None:
        """Simulate disconnecting from an in-memory database."""
        logging.info("Disconnected from in-memory database")

    def register(self, incident: Incident) -> None:
        """Register a new incident in the database system."""
        new_db_entry = {incident.id: incident}
        self._storage.update(new_db_entry)

    def update_entry(self, incident: Incident) -> None:
        """Update an existing incident entry in the database."""
        if incident.id not in self._storage:
            raise IncidentNotFoundError(incident.id)
        self._storage[incident.id] = incident

    def get_by_id(self, incident_id: str) -> Incident:
        """Retrieve an incident by its ID."""
        incident_by_id = self._storage.get(incident_id)
        if incident_by_id is None:
            raise IncidentNotFoundError(incident_id)
        return incident_by_id

    def get_all(self) -> dict[str, Incident]:
        """Retrieve all incidents from the database."""
        return self._storage

    def remove_by_id(self, incident_id: str) -> None:
        """Remove an incident from the database by its ID."""
        if incident_id not in self._storage:
            raise IncidentNotFoundError(incident_id)
        del self._storage[incident_id]


_db_singleton = InMemoryClient()


def get_database_client() -> DataClient:
    """Exposed factory function to get the appropriate database client."""
    return _db_singleton
