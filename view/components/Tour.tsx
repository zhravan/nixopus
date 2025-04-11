'use client';
import { TourProvider, StepType, useTour as useReactourTour } from '@reactour/tour';
import { useTour as useCustomTour } from '../hooks/useTour';
import React from 'react';

const steps: StepType[] = [
  {
    selector: '[data-slot="sidebar-trigger"]',
    content: 'This is the sidebar toggle. Click it to show/hide the navigation menu.'
  },
  {
    selector: '[data-sidebar="sidebar"]',
    content:
      'This is your main navigation menu. Here you can access all the different sections of Nixopus.'
  },
  {
    selector: '[data-slot="keyboard-shortcuts"]',
    content: 'Click here to view all available keyboard shortcuts for quick navigation.'
  },
  {
    selector: '[data-slot="terminal"]',
    content: 'This is the terminal. You can use it to execute commands and manage your environment.'
  },
  {
    selector: '[data-slot="breadcrumb"]',
    content: 'These breadcrumbs show your current location in the application.'
  },
  {
    selector: '[data-slot="tour-trigger"]',
    content: 'Click this button anytime to restart the tour and learn more about Nixopus.'
  }
];

const TourContent = ({ children }: { children: React.ReactNode }) => {
  const { setIsOpen } = useReactourTour();
  const { hasSeenTour, startTour, stopTour } = useCustomTour();

  React.useEffect(() => {
    if (!hasSeenTour) {
      setTimeout(() => {
        setIsOpen(true);
      }, 1000);
    }
  }, [hasSeenTour, setIsOpen]);

  React.useEffect(() => {
    const handleTourTrigger = () => {
      startTour();
      setIsOpen(true);
    };

    const tourTrigger = document.querySelector('[data-slot="tour-trigger"]');
    if (tourTrigger) {
      tourTrigger.addEventListener('click', handleTourTrigger);
      return () => tourTrigger.removeEventListener('click', handleTourTrigger);
    }
  }, [startTour, setIsOpen]);

  return <>{children}</>;
};

export const Tour = ({ children }: { children: React.ReactNode }) => {
  const { stopTour } = useCustomTour();

  return (
    <TourProvider
      steps={steps}
      showNavigation={true}
      showBadge={false}
      showDots={true}
      showCloseButton={true}
      disableInteraction={false}
      disableKeyboardNavigation={false}
      disableDotsNavigation={false}
      disableFocusLock={false}
      beforeClose={() => {
        stopTour();
        return true;
      }}
      styles={{
        popover: (base: any) => ({
          ...base,
          borderRadius: 'var(--radius)',
          backgroundColor: 'var(--secondary)',
          color: 'var(--secondary-foreground)'
        }),
        dot: (base: any, state: any) => ({
          ...base,
          background: state?.current ? 'var(--primary)' : 'var(--muted)'
        }),
        navigation: (base: any) => ({
          ...base,
          display: 'flex',
          gap: '0.5rem'
        }),
        close: (base: any) => ({
          ...base,
          color: 'var(--secondary-foreground)',
          cursor: 'pointer',
          '&:hover': {
            opacity: 0.8
          }
        })
      }}
    >
      <TourContent>{children}</TourContent>
    </TourProvider>
  );
};
