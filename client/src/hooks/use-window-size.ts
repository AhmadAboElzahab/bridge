import * as React from "react";
import { useIsomorphicLayoutEffect } from "@/hooks/use-isomorphic-layout-effect";

interface WindowSize {
  width: number;
  height: number;
}

interface UseWindowSizeProps {
  defaultWidth?: number;
  defaultHeight?: number;
}

export function useWindowSize(props: UseWindowSizeProps = {}): WindowSize {
  const { defaultWidth = 0, defaultHeight = 0 } = props;

  const [windowSize, setWindowSize] = React.useState<WindowSize>(() => {
    if (typeof window !== "undefined") {
      return { width: window.innerWidth, height: window.innerHeight };
    }
    return { width: defaultWidth, height: defaultHeight };
  });

  useIsomorphicLayoutEffect(() => {
    let timeoutId: NodeJS.Timeout | null = null;

    function onResize() {
      if (timeoutId) {
        clearTimeout(timeoutId);
      }

      timeoutId = setTimeout(() => {
        setWindowSize({
          width: window.innerWidth,
          height: window.innerHeight,
        });
      }, 150);
    }

    window.addEventListener("resize", onResize);
    return () => {
      window.removeEventListener("resize", onResize);
      if (timeoutId) {
        clearTimeout(timeoutId);
      }
    };
  }, []);

  return windowSize;
}
