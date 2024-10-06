"use client"

import i18n from "@/config/i18n"

export default function Submit({ children, isSubmitting }: { children: React.ReactNode; isSubmitting: boolean }) {
  return (
    <>
      {isSubmitting ? (
        <button aria-busy="true">{i18n.t("common.loading")}</button>
      ) : (
        <button type="submit" disabled={isSubmitting}>
          {children}
        </button>
      )}
    </>
  )
}
