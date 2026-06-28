export function addHeadingIds(html) {
  if (!html) return html
  const parser = new DOMParser()
  const doc = parser.parseFromString(html, 'text/html')
  doc.querySelectorAll('h1, h2, h3, h4, h5, h6').forEach(h => {
    const text = h.textContent.trim()
    if (!h.id && text) {
      h.id = text
        .toLowerCase()
        .replace(/[^\p{L}\p{N}\s-]/gu, '')
        .replace(/\s+/g, '-')
        .replace(/-+/g, '-')
        .replace(/^-|-$/g, '')
    }
  })
  return doc.body.innerHTML
}
