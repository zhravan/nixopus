'use client';
import React, { useEffect, useState, useRef } from 'react';
import '@xterm/xterm/css/xterm.css';
import { useTerminal } from './utils/useTerminal';
import { useContainerReady } from './utils/isContainerReady';
import { X } from 'lucide-react';

type TerminalProps = {
  isOpen: boolean;
  toggleTerminal: () => void;
  isTerminalOpen: boolean;
  setFitAddonRef: React.Dispatch<React.SetStateAction<any | null>>;
};

export const Terminal: React.FC<TerminalProps> = ({
  isOpen,
  toggleTerminal,
  isTerminalOpen,
  setFitAddonRef
}) => {
  const [dimensions, setDimensions] = useState({ width: 0, height: 0 });
  const containerRef = useRef<HTMLDivElement>(null);

  const { terminalRef, fitAddonRef, initializeTerminal, destroyTerminal } = useTerminal(
    isTerminalOpen,
    dimensions.width,
    dimensions.height
  ) as {
    terminalRef: React.RefObject<HTMLDivElement>;
    fitAddonRef: any;
    initializeTerminal: () => void;
    destroyTerminal: () => void;
  };

  const isContainerReady = useContainerReady(isTerminalOpen, terminalRef);

  useEffect(() => {
    if (!containerRef.current) return;

    const updateDimensions = () => {
      if (containerRef.current) {
        setDimensions({
          width: containerRef.current.offsetWidth,
          height: containerRef.current.offsetHeight
        });
      }
    };

    updateDimensions();

    const resizeObserver = new ResizeObserver(updateDimensions);
    resizeObserver.observe(containerRef.current);

    return () => {
      resizeObserver.disconnect();
    };
  }, [isTerminalOpen]);

  useEffect(() => {
    if (isTerminalOpen && isContainerReady) {
      setTimeout(initializeTerminal, 0);
    } else {
      destroyTerminal();
    }
  }, [isTerminalOpen, isContainerReady, initializeTerminal, destroyTerminal]);

  useEffect(() => {
    if (fitAddonRef) {
      setFitAddonRef(fitAddonRef);
    }
  }, [fitAddonRef, setFitAddonRef]);

  return (
    <div className="flex h-full flex-col" ref={containerRef}>
      <div className="flex h-5 items-center justify-between bg-secondary px-1 py-2 opacity-50">
        <span className="text-xs">
          Terminal <span className="text-xs">⌘</span>J
        </span>
        <X className="h-4 w-4 hover:text-destructive" onClick={toggleTerminal} />
      </div>
      <div
        ref={terminalRef}
        className="flex-grow overflow-hidden bg-secondary"
        style={{
          height: isTerminalOpen ? 'calc(100% - 32px)' : '0',
          visibility: isTerminalOpen ? 'visible' : 'hidden'
        }}
      />
    </div>
  );
};
