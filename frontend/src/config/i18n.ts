import i18n from "i18next"
import { initReactI18next } from "react-i18next"

import en from "../../public/locales/en/translation.json"

i18n.use(initReactI18next).init({
  compatibilityJSON: "v3",
  resources: { en },
  lng: "en",
  fallbackLng: "en",
  debug: false,
  interpolation: {
    escapeValue: false,
  },
  react: {
    useSuspense: false,
  },
  returnNull: false,
  returnObjects: true,
})

export default i18n
