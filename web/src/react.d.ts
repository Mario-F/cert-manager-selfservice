// Workaround: https://github.com/johnsoncodehk/volar/discussions/592
import { AriaAttributes, DOMAttributes } from 'react' // not needed if skipLibCheck = true

declare module 'react' {
  interface HTMLAttributes<T> extends AriaAttributes, DOMAttributes<T> {
    // extends React's HTMLAttributes
    class?: string
  }
}
