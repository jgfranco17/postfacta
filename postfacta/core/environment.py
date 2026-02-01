import logging
import os
from enum import StrEnum


class EnvVariable(StrEnum):
    LOG_LEVEL = "LOG_LEVEL"


class LogLevel(StrEnum):
    CRITICAL = "CRITICAL"
    ERROR = "ERROR"
    WARNING = "WARNING"
    INFO = "INFO"
    DEBUG = "DEBUG"
    NOTSET = "NOTSET"


def get_logging_level_from_env() -> int:
    """Get logging level from environment variable."""
    level_str = os.getenv(EnvVariable.LOG_LEVEL, LogLevel.DEBUG)
    levels_mapping = {
        LogLevel.CRITICAL.value: logging.CRITICAL,
        LogLevel.ERROR.value: logging.ERROR,
        LogLevel.WARNING.value: logging.WARNING,
        LogLevel.INFO.value: logging.INFO,
        LogLevel.DEBUG.value: logging.DEBUG,
        LogLevel.NOTSET.value: logging.NOTSET,
    }
    return levels_mapping.get(level_str, logging.INFO)
