"use client"

import { Field, FieldProps } from "formik"

import styles from "./Input.module.css"

interface InputProps {
  type: string
  name: string
  placeholder: string
  autoComplete?: string
  label?: string
}

interface FieldMeta {
  touched: boolean
  error?: string
}

export default function Input(props: InputProps) {
  const isInvalid = (meta: FieldMeta) => meta.touched && !!meta.error

  return (
    <Field {...props}>
      {({ field, meta }: FieldProps<string> & { meta: FieldMeta }) => (
        <>
          <input
            {...field}
            type={props.type}
            placeholder={props.placeholder}
            autoComplete={props.autoComplete}
            aria-invalid={isInvalid(meta) ? "true" : undefined}
            aria-describedby={isInvalid(meta) ? `${field.name}-helper` : undefined}
            className={isInvalid(meta) ? styles.invalid : ""}
          />
          {isInvalid(meta) && (
            <small className={styles.small} id={`${field.name}-helper`}>
              {meta.error}
            </small>
          )}
        </>
      )}
    </Field>
  )
}
