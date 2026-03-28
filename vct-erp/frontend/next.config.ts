import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  experimental: {
    reactCompiler: process.env.NODE_ENV === "production",
  },
};

export default nextConfig;
