"use client"

import { useState } from "react"
import { Formik, Form } from "formik"
import Image from "next/image"
import Link from "next/link"

import i18n from "@/config/i18n"
import { API_URL } from "@/config/env"
import { URL_REGEX } from "@/helpers/url"
import Button from "@/components/ui/Button"
import Fieldset from "@/components/ui/Fieldset"
import Input from "@/components/ui/Input"
import Submit from "@/components/ui/Submit"
import styles from "./page.module.css"

export default function Home() {
  const [shortUrl, setShortUrl] = useState("")

  const initialValues = {
    url: "",
  }

  const handleValidate = (values: { url: string }) => {
    const errors = {} as { url?: string }
    if (!values.url) {
      errors.url = i18n.t("validation.field.required", { field: i18n.t("url.title") })
    } else if (!URL_REGEX.test(values.url)) {
      errors.url = i18n.t("validation.field.invalid", { field: i18n.t("url.title") })
    }

    return errors
  }

  const handleSubmit = async (values: { url: string }, setSubmitting: (isSubmitting: boolean) => void) => {
    if (API_URL === undefined) {
      console.error("--- error ---")
      console.error("API_URL is not defined")
      return
    }

    await fetch(API_URL, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(values.url),
    })
      .then((response) => {
        response.text().then((shortUrl) => {
          setShortUrl(shortUrl)
        })
      })
      .catch((error) => {
        console.error("--- error ---")
        console.error(error)
      })
      .finally(() => {
        setSubmitting(false)
      })
  }

  const handleReset = () => {
    setShortUrl("")
  }

  return (
    <main className={styles.main}>
      <div className={styles.container}>
        {shortUrl ? (
          <div className={styles.content}>
            <Link className={styles.link} href={shortUrl}>
              {shortUrl}
            </Link>
            <Button className={styles.back} onClick={handleReset}>
              {i18n.t("common.back")}
            </Button>
          </div>
        ) : (
          <Formik
            className={styles.container}
            initialValues={initialValues}
            validate={(values) => handleValidate(values)}
            onSubmit={(values, { setSubmitting }) => handleSubmit(values, setSubmitting)}
          >
            {({ isSubmitting }) => (
              <Form>
                <Fieldset>
                  <Input type="text" name="url" placeholder={i18n.t("url.form.placeholder")} autoComplete="off" />
                  <Submit isSubmitting={isSubmitting}>{i18n.t("common.submit")}</Submit>
                </Fieldset>
              </Form>
            )}
          </Formik>
        )}
      </div>
      <footer className={styles.footer}>
        <a
          href="https://nextjs.org/learn?utm_source=create-next-app&utm_medium=appdir-template&utm_campaign=create-next-app"
          target="_blank"
          rel="noopener noreferrer"
        >
          <Image aria-hidden src="https://nextjs.org/icons/file.svg" alt="File icon" width={16} height={16} />
          Learn
        </a>
        <a
          href="https://vercel.com/templates?framework=next.js&utm_source=create-next-app&utm_medium=appdir-template&utm_campaign=create-next-app"
          target="_blank"
          rel="noopener noreferrer"
        >
          <Image aria-hidden src="https://nextjs.org/icons/window.svg" alt="Window icon" width={16} height={16} />
          Examples
        </a>
        <a
          href="https://nextjs.org?utm_source=create-next-app&utm_medium=appdir-template&utm_campaign=create-next-app"
          target="_blank"
          rel="noopener noreferrer"
        >
          <Image aria-hidden src="https://nextjs.org/icons/globe.svg" alt="Globe icon" width={16} height={16} />
          Go to nextjs.org →
        </a>
      </footer>
    </main>
  )
}
