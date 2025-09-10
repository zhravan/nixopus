import json

import typer

from app.utils.config import DEFAULT_COMPOSE_FILE, NIXOPUS_CONFIG_DIR, Config
from app.utils.logger import Logger
from app.utils.output_formatter import OutputFormatter
from app.utils.timeout import TimeoutWrapper

from .down import Down, DownConfig
from .messages import (
    services_restarted_successfully,
    services_started_successfully,
    services_status_retrieved,
    services_stopped_successfully,
)
from .ps import Ps, PsConfig
from .restart import Restart, RestartConfig
from .up import Up, UpConfig

service_app = typer.Typer(help="Manage Nixopus services")

config = Config()
nixopus_config_dir = config.get_yaml_value(NIXOPUS_CONFIG_DIR)
compose_file = config.get_yaml_value(DEFAULT_COMPOSE_FILE)
compose_file_path = nixopus_config_dir + "/" + compose_file


@service_app.command()
def up(
    name: str = typer.Option("all", "--name", "-n", help="The name of the service to start, defaults to all"),
    verbose: bool = typer.Option(False, "--verbose", "-v", help="Verbose output"),
    output: str = typer.Option("text", "--output", "-o", help="Output format, text, json"),
    dry_run: bool = typer.Option(False, "--dry-run", help="Dry run"),
    detach: bool = typer.Option(False, "--detach", "-d", help="Detach from the service and run in the background"),
    env_file: str = typer.Option(None, "--env-file", "-e", help="Path to the environment file"),
    compose_file: str = typer.Option(compose_file_path, "--compose-file", "-f", help="Path to the compose file"),
    timeout: int = typer.Option(10, "--timeout", "-t", help="Timeout in seconds"),
):
    """Start Nixopus services"""
    logger = Logger(verbose=verbose)

    try:
        config = UpConfig(
            name=name,
            detach=detach,
            env_file=env_file,
            verbose=verbose,
            output=output,
            dry_run=dry_run,
            compose_file=compose_file,
        )

        up_service = Up(logger=logger)

        with TimeoutWrapper(timeout):
            if config.dry_run:
                formatted_output = up_service.format_dry_run(config)
                logger.info(formatted_output)
                return
            else:
                result = up_service.up(config)

        if result.success:
            formatted_output = up_service.format_output(result, output)
            if output == "json":
                logger.info(formatted_output)
            else:
                logger.success(services_started_successfully.format(services=result.name))
                if formatted_output:
                    logger.info(formatted_output)
        else:
            logger.error(result.error if result.error is not None else "Unknown error")
            raise typer.Exit(1)

    except TimeoutError as e:
        logger.error(e)
        raise typer.Exit(1)
    except Exception as e:
        logger.error(str(e))
        raise typer.Exit(1)


@service_app.command()
def down(
    name: str = typer.Option("all", "--name", "-n", help="The name of the service to stop, defaults to all"),
    verbose: bool = typer.Option(False, "--verbose", "-v", help="Verbose output"),
    output: str = typer.Option("text", "--output", "-o", help="Output format, text, json"),
    dry_run: bool = typer.Option(False, "--dry-run", help="Dry run"),
    env_file: str = typer.Option(None, "--env-file", "-e", help="Path to the environment file"),
    compose_file: str = typer.Option(compose_file_path, "--compose-file", "-f", help="Path to the compose file"),
    timeout: int = typer.Option(10, "--timeout", "-t", help="Timeout in seconds"),
):
    """Stop Nixopus services"""
    logger = Logger(verbose=verbose)

    try:
        config = DownConfig(
            name=name, env_file=env_file, verbose=verbose, output=output, dry_run=dry_run, compose_file=compose_file
        )

        down_service = Down(logger=logger)

        with TimeoutWrapper(timeout):
            if config.dry_run:
                formatted_output = down_service.format_dry_run(config)
                logger.info(formatted_output)
                return
            else:
                result = down_service.down(config)

        if result.success:
            formatted_output = down_service.format_output(result, output)
            if output == "json":
                logger.info(formatted_output)
            else:
                logger.success(services_stopped_successfully.format(services=result.name))
                if formatted_output:
                    logger.info(formatted_output)
        else:
            logger.error(result.error)
            raise typer.Exit(1)

    except TimeoutError as e:
        logger.error(e)
        raise typer.Exit(1)
    except Exception as e:
        logger.error(str(e))
        raise typer.Exit(1)


@service_app.command()
def ps(
    name: str = typer.Option("all", "--name", "-n", help="The name of the service to show, defaults to all"),
    verbose: bool = typer.Option(False, "--verbose", "-v", help="Verbose output"),
    output: str = typer.Option("text", "--output", "-o", help="Output format, text, json"),
    dry_run: bool = typer.Option(False, "--dry-run", "-d", help="Dry run"),
    env_file: str = typer.Option(None, "--env-file", "-e", help="Path to the environment file"),
    compose_file: str = typer.Option(compose_file_path, "--compose-file", "-f", help="Path to the compose file"),
    timeout: int = typer.Option(10, "--timeout", "-t", help="Timeout in seconds"),
):
    """Show status of Nixopus services"""
    logger = Logger(verbose=verbose)

    try:
        config = PsConfig(
            name=name, env_file=env_file, verbose=verbose, output=output, dry_run=dry_run, compose_file=compose_file
        )

        ps_service = Ps(logger=logger)

        with TimeoutWrapper(timeout):
            if config.dry_run:
                formatted_output = ps_service.format_dry_run(config)
                logger.info(formatted_output)
                return
            else:
                result = ps_service.ps(config)

        if result.success:
            formatted_output = ps_service.format_output(result, output)
            logger.info(formatted_output)
        else:
            logger.error(result.error)
            raise typer.Exit(1)

    except TimeoutError as e:
        logger.error(e)
        raise typer.Exit(1)
    except Exception as e:
        logger.error(str(e))
        raise typer.Exit(1)


@service_app.command()
def restart(
    name: str = typer.Option("all", "--name", "-n", help="The name of the service to restart, defaults to all"),
    verbose: bool = typer.Option(False, "--verbose", "-v", help="Verbose output"),
    output: str = typer.Option("text", "--output", "-o", help="Output format, text, json"),
    dry_run: bool = typer.Option(False, "--dry-run", "-d", help="Dry run"),
    env_file: str = typer.Option(None, "--env-file", "-e", help="Path to the environment file"),
    compose_file: str = typer.Option(compose_file_path, "--compose-file", "-f", help="Path to the compose file"),
    timeout: int = typer.Option(10, "--timeout", "-t", help="Timeout in seconds"),
):
    """Restart Nixopus services"""
    logger = Logger(verbose=verbose)

    try:
        config = RestartConfig(
            name=name, env_file=env_file, verbose=verbose, output=output, dry_run=dry_run, compose_file=compose_file
        )

        restart_service = Restart(logger=logger)

        with TimeoutWrapper(timeout):
            if config.dry_run:
                formatted_output = restart_service.format_dry_run(config)
                logger.info(formatted_output)
                return
            else:
                result = restart_service.restart(config)

        if result.success:
            formatted_output = restart_service.format_output(result, output)
            if output == "json":
                logger.info(formatted_output)
            else:
                logger.success(services_restarted_successfully.format(services=result.name))
                if formatted_output:
                    logger.info(formatted_output)
        else:
            logger.error(result.error)
            raise typer.Exit(1)

    except TimeoutError as e:
        logger.error(e)
        raise typer.Exit(1)
    except Exception as e:
        logger.error(str(e))
        raise typer.Exit(1)
