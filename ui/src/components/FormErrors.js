export default function FormErrors({ errors }) {
  if (errors.length === 0) {
    return null;
  }

  const errorElements = errors.map((error, i) => {
    return(<li key={`error-${i}`}>{error}</li>);
  })

  return(
    <blockquote>
      <ul>
        {errorElements}
      </ul>
    </blockquote>
  )
}
