import bundleAnalyzer from "@next/bundle-analyzer";
import type { NextConfig } from "next";

import "./src/env.js";

const withBundleAnalyzer = bundleAnalyzer({
  enabled: process.env.ANALYZE === "true",
});

const nextConfig: NextConfig = {
  experimental: {
    optimizePackageImports: [
      "@tanstack/react-query",
      "@tanstack/react-table",
      "@tanstack/react-virtual",
    ],
  },
  async rewrites() {
    const apiBase = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";
    return [
      {
        source: "/api/:path*",
        destination: `${apiBase}/api/:path*`,
      },
      {
        source: "/storage/:path*",
        destination: `${apiBase}/storage/:path*`,
      },
    ];
  },
};

export default withBundleAnalyzer(nextConfig);
