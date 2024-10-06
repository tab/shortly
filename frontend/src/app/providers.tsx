"use client"

import { ThemeProvider } from "next-themes"

const Providers = ({
  children,
}: Readonly<{
  children: React.ReactNode
}>) => {
  return <ThemeProvider>{children}</ThemeProvider>
}

export default Providers
