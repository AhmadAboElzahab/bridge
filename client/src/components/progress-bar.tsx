"use client";

import NProgress from "nprogress";
import { usePathname, useSearchParams } from "next/navigation";
import { useEffect, useRef } from "react";
import { useIsFetching } from "@tanstack/react-query";

NProgress.configure({ showSpinner: false, trickleSpeed: 200 });

function RouteProgressBar() {
  const pathname = usePathname();
  const searchParams = useSearchParams();
  const mounted = useRef(false);

  useEffect(() => {
    if (!mounted.current) {
      mounted.current = true;
      return;
    }
    NProgress.done();
  }, [pathname, searchParams]);

  useEffect(() => {
    const handleClick = (e: MouseEvent) => {
      const anchor = (e.target as HTMLElement).closest("a");
      if (!anchor) return;
      const href = anchor.getAttribute("href");
      if (!href || href.startsWith("#") || href.startsWith("mailto:")) return;
      const url = new URL(href, window.location.href);
      if (url.origin !== window.location.origin) return;
      if (url.pathname !== window.location.pathname || url.search !== window.location.search) {
        NProgress.start();
      }
    };
    document.addEventListener("click", handleClick);
    return () => document.removeEventListener("click", handleClick);
  }, []);

  return null;
}

function QueryProgressBar() {
  const isFetching = useIsFetching();

  useEffect(() => {
    if (isFetching > 0) {
      NProgress.start();
    } else {
      NProgress.done();
    }
  }, [isFetching]);

  return null;
}

export function ProgressBar() {
  return (
    <>
      <RouteProgressBar />
      <QueryProgressBar />
    </>
  );
}
