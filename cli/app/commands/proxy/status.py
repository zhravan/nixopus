from typing import Protocol

from app.utils.config import PROXY_PORT, Config
from app.utils.protocols import LoggerProtocol

from .base import BaseAction, BaseCaddyCommandBuilder, BaseCaddyService, BaseConfig, BaseFormatter, BaseResult, BaseService
from .messages import (
    dry_run_command,
    dry_run_command_would_be_executed,
    dry_run_mode,
    dry_run_port,
    end_dry_run,
)

config = Config()
proxy_port = config.get_yaml_value(PROXY_PORT)


class CaddyServiceProtocol(Protocol):
    def check_status(self, port: int = proxy_port) -> tuple[bool, str]: ...


class CaddyCommandBuilder(BaseCaddyCommandBuilder):
    @staticmethod
    def build_status_command(port: int = proxy_port) -> list[str]:
        return BaseCaddyCommandBuilder.build_status_command(port)


class StatusFormatter(BaseFormatter):
    def format_output(self, result: "StatusResult", output: str) -> str:
        if output == "json":
            status_msg = "Caddy is running" if result.success else (result.error or "Caddy not running")
            return super().format_output(result, output, status_msg, result.error or "Caddy not running")

        if result.success:
            return "Caddy is running"
        else:
            return result.error or "Caddy not running"

    def format_dry_run(self, config: "StatusConfig") -> str:
        dry_run_messages = {
            "mode": dry_run_mode,
            "command_would_be_executed": dry_run_command_would_be_executed,
            "command": dry_run_command,
            "port": dry_run_port,
            "end": end_dry_run,
        }
        return super().format_dry_run(config, CaddyCommandBuilder(), dry_run_messages)


class CaddyService(BaseCaddyService):
    def __init__(self, logger: LoggerProtocol):
        super().__init__(logger)

    def get_status(self, port: int = proxy_port) -> tuple[bool, str]:
        return self.check_status(port)


class StatusResult(BaseResult):
    pass


class StatusConfig(BaseConfig):
    pass


class StatusService(BaseService[StatusConfig, StatusResult]):
    def __init__(self, config: StatusConfig, logger: LoggerProtocol = None, caddy_service: CaddyServiceProtocol = None):
        super().__init__(config, logger, caddy_service)
        self.caddy_service = caddy_service or CaddyService(self.logger)
        self.formatter = StatusFormatter()

    def _create_result(self, success: bool, error: str = None) -> StatusResult:
        return StatusResult(
            proxy_port=self.config.proxy_port,
            verbose=self.config.verbose,
            output=self.config.output,
            success=success,
            error=error,
        )

    def status(self) -> StatusResult:
        return self.execute()

    def execute(self) -> StatusResult:
        success, message = self.caddy_service.get_status(self.config.proxy_port)
        return self._create_result(success, None if success else message)

    def status_and_format(self) -> str:
        return self.execute_and_format()

    def execute_and_format(self) -> str:
        if self.config.dry_run:
            return self.formatter.format_dry_run(self.config)

        result = self.execute()
        return self.formatter.format_output(result, self.config.output)


class Status(BaseAction[StatusConfig, StatusResult]):
    def __init__(self, logger: LoggerProtocol = None):
        super().__init__(logger)
        self.formatter = StatusFormatter()

    def status(self, config: StatusConfig) -> StatusResult:
        return self.execute(config)

    def execute(self, config: StatusConfig) -> StatusResult:
        service = StatusService(config, logger=self.logger)
        return service.execute()

    def format_output(self, result: StatusResult, output: str) -> str:
        return self.formatter.format_output(result, output)
