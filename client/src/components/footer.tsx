import * as Constants from '@/constants';
export default function Footer() {
  return (
    <footer className="m500:text-sm dark:bg-secondaryBlack z-30 bg-white px-5 py-5 text-center font-base">
      {/* Released under MIT License. The source code is available on{' '} */}
      Entre em contato via{' '}
      <a
        // target="_blank"
        href={"mailto:" + Constants.CONTACT_EMAIL}
        className="font-heading underline"
      >
        {Constants.CONTACT_EMAIL}
      </a>
      .
    </footer>
  )
}
