/**
 * Converts Base64 encoded string of an image 
 * to a BLOB and returns the URL of the image BLOB
 */
export function convertBase64ToImage(base64Data: string): string {
    const binaryStr = atob(base64Data)

    const byteArray = new Uint8Array(binaryStr.length)
    for (let i = 0; i < binaryStr.length; i++) {
        byteArray[i] = binaryStr.charCodeAt(i)
    }

    const blob = new Blob([byteArray], { type: "image/png" })
    const url = URL.createObjectURL(blob)

    return url
}