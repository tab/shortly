"use client"

import { useEffect, useState } from "react"
import { API_URL } from "@/config/env"

export default function Page({ params }: { params: { id: string } }) {
  const [isFetching, setIsFetching] = useState(false)
  const [longUrl, setLongUrl] = useState<string | null>(null)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (isFetching) {
      return
    }

    handleFetch()
  }, [params.id])

  const handleFetch = async () => {
    setIsFetching(true)
    setError(null)

    if (!API_URL) {
      console.error("--- error ---")
      console.error("API_URL is not defined")
      setIsFetching(false)
      return
    }

    try {
      const response = await fetch(`${API_URL}/${params.id}`, {
        method: "GET",
        headers: {
          "Content-Type": "text/plain",
        },
      })

      if (response.status === 307) {
        // Manually handle redirect
        const redirectUrl = response.headers.get("Location")
        if (redirectUrl) {
          // Optionally update state or handle it here
          window.location.href = redirectUrl // Redirect the browser
        }
      } else if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      } else {
        const longUrl = await response.text()
        setLongUrl(longUrl)
      }
    } catch (error) {
      console.error("--- error ---")
      console.error(error)
      setError("Failed to fetch the long URL.")
    } finally {
      setIsFetching(false)
    }
  }

  return (
    <div>
      <h1>Code: {params.id}</h1>
      {isFetching ? <p>Loading...</p> : error ? <p>{error}</p> : longUrl ? <p>Redirecting to: {longUrl}</p> : <p>No URL found.</p>}
    </div>
  )
}
