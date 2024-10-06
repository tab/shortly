"use client"

export default function Fieldset({ children }: { children: React.ReactNode }) {
  return <fieldset role="group">{children}</fieldset>
}
