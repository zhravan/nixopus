import { useState, useEffect } from 'react';

const TOUR_SEEN_KEY = 'nixopus_tour_seen';

export const useTour = () => {
  const [hasSeenTour, setHasSeenTour] = useState(false);

  useEffect(() => {
    const seen = localStorage.getItem(TOUR_SEEN_KEY);
    setHasSeenTour(!!seen);
  }, []);

  const startTour = () => {
    setHasSeenTour(false);
  };

  const stopTour = () => {
    setHasSeenTour(true);
    localStorage.setItem(TOUR_SEEN_KEY, 'true');
  };

  return {
    hasSeenTour,
    startTour,
    stopTour
  };
};
