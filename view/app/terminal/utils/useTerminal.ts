import { useState, useRef, useCallback, useEffect } from 'react';
import { StopExecution } from './stopExecution';
import { useTheme } from 'next-themes';
import {
    handleBackspaceKey,
    handleControlKey,
    handleEnterKey,
    handlePrintableKey,
    handleTabKey,
    specialKeyCodeMappings,
} from './utils';

// what are the type of output that possible from the backend
enum outputType {
    STDOUT = 'stdout',
    STDERR = 'stderr',
    EXIT = 'exit',
}

// this is the response type that we get from the backend
type TerminalOutput = {
    data: {
        output_type: string;
        content: string;
    };
};

export const useTerminal = () => {
    const terminalRef = useRef<HTMLDivElement | null>(null);
    const fitAddonRef = useRef<any | null>(null);
    const [terminalInstance, setTerminalInstance] = useState<any | null>(null);
    const { isStopped, setIsStopped } = StopExecution();
    const { resolvedTheme } = useTheme();

    // // websocket hooks and initialization
    // const { lastMessage, sendJsonMessage } = useWebSocket(WS_URL, {
    //     shouldReconnect: () => true,
    //     onOpen: () => console.log('WebSocket connection established.'),
    //     onClose: () => console.log('WebSocket connection closed.'),
    //     onError: (event) => console.error('WebSocket error:', event),
    // });

    // this will be triggered to stop the execution of the terminal
    useEffect(() => {
        if (isStopped && terminalInstance) {
            console.log('Sending SIGINT');
            // sendJsonMessage({
            //     TerminalCommand: {
            //         command: '\x03',
            //     },
            // });
            setIsStopped(false);
        }
    }, [isStopped, terminalInstance, setIsStopped]);

    // // this will be triggered whenever there is a new message and it will update the terminal
    // useEffect(() => {
    //     if (lastMessage !== null && terminalInstance) {
    //         try {
    //             const parsedMessage: TerminalOutput = JSON.parse(lastMessage.data);
    //             if (parsedMessage.data.output_type === outputType.EXIT) {
    //                 setTerminalInstance(null);
    //                 return;
    //             }
    //             if (parsedMessage.data.output_type === outputType.STDERR) {
    //                 terminalInstance.writeln('\x1B[31m' + parsedMessage.data.content + '\x1B[0m');
    //                 return;
    //             }

    //             const output = parsedMessage.data.content;
    //             console.log(output);
    //             terminalInstance.write(output);
    //         } catch (error) {
    //             console.error('Error parsing WebSocket message:', error);
    //             if (terminalInstance) {
    //                 terminalInstance.writeln('\r\nError: Failed to parse server response');
    //             }
    //         }
    //     }
    // }, [lastMessage, terminalInstance]);

    // this will be triggered whenever the theme changes
    useEffect(() => {
        if (resolvedTheme && terminalInstance) {
            const newColors = {
                foreground: "red",
                background: "green",
                cursor: "red",
            };
            terminalInstance.options.theme = { ...newColors };
        }
    }, [resolvedTheme]);

    const initializeTerminal = useCallback(async () => {
        if (!terminalRef.current || terminalInstance) return;
        // lazy load xterm and related plugins
        const { Terminal } = await import('@xterm/xterm');
        const { FitAddon } = await import('xterm-addon-fit');
        const { WebLinksAddon } = await import('xterm-addon-web-links');

        // initialize xterm
        const term = new Terminal({
            cursorBlink: true,
            fontFamily: '"Menlo", "DejaVu Sans Mono", "Consolas", monospace',
            fontSize: 14,
            theme: {
                foreground: "red",
                background: `black`,
                cursor: "red",
            },
            allowTransparency: true,
            rightClickSelectsWord: true,
            // logLevel: 'debug',
        });

        // initialize addons
        const fitAddon = new FitAddon();
        const webLinksAddon = new WebLinksAddon();

        // load addons and activate it
        term.loadAddon(fitAddon);
        term.loadAddon(webLinksAddon);
        fitAddonRef.current = fitAddon;
        term.open(terminalRef.current);
        fitAddon.activate(term);
        fitAddon.fit();

        // set padding for the terminal
        if (terminalRef.current) {
            terminalRef.current.style.padding = '10px';
        }
        let currentLine = '';

        // whenever there is something written onto the terminal this will be called and updates the currentline
        term.onData((data) => {
            currentLine += data;
        });

        term.onKey(({ key, domEvent }) => {
            const printable = !domEvent.altKey && !domEvent.ctrlKey && !domEvent.metaKey;

            if (domEvent.keyCode === 13) {
                handleEnterKey(term, currentLine, (data: any) => {
                    
                });
                return;
            } else if (domEvent.keyCode === 9) {
                handleTabKey(currentLine, (data: any) => {}, term);
                return;
            } else if (domEvent.keyCode === 8) {
                currentLine = handleBackspaceKey(term, currentLine);
                return;
            }
            // these keys are for special keys, like arrow up down left right home delete escape
            else if (
                domEvent.keyCode === 38 ||
                domEvent.keyCode === 40 ||
                domEvent.keyCode === 37 ||
                domEvent.keyCode === 39
            ) {
                // sendJsonMessage({
                //     TerminalCommand: {
                //         command: specialKeyCodeMappings[domEvent.keyCode],
                //     },
                // });
            } else if (printable) {
                handlePrintableKey(term, currentLine, key);
            } else {
                let command = '';
                if (domEvent.ctrlKey) {
                    command = handleControlKey(domEvent.key.toLowerCase());
                }
                if (command) {
                    // sendJsonMessage({
                    //     TerminalCommand: {
                    //         command: command,
                    //     },
                    // });
                }
            }
        });

        setTerminalInstance(term);
    }, [terminalInstance]);

    // clean up terminal
    const destroyTerminal = useCallback(() => {
        if (terminalInstance) {
            terminalInstance.dispose();
            setTerminalInstance(null);
        }
    }, [terminalInstance]);

    // this will be triggered for destroying the terminal
    useEffect(() => {
        return () => {
            destroyTerminal();
        };
    }, [destroyTerminal]);

    return { terminalRef, fitAddonRef, initializeTerminal, destroyTerminal, terminalInstance };
};
