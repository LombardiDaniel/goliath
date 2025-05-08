import { ThemeProvider } from '@/components/theme-provider';
import { ThemeSwitcher } from '@/components/theme-switcher';
import * as Constants from '@/constants';
import type { Metadata } from 'next';
import { DM_Sans } from 'next/font/google';
import './globals.css';

const dmSans = DM_Sans({ subsets: ['latin'] })

export const metadata: Metadata = {
  metadataBase: new URL(Constants.URL),
  title: Constants.APP_TITLE,
  description: Constants.DESCRIPTION[0],
  openGraph: {
    title: Constants.APP_TITLE,
    description: Constants.DESCRIPTION[0],
    url: Constants.URL,
    siteName: Constants.APP_TITLE,
  },
}

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode
}>) {
  return (
    <html lang="en">
      <body className={dmSans.className}>
        <ThemeProvider attribute="class" disableTransitionOnChange>
          {children}
          <ThemeSwitcher />
        </ThemeProvider>
      </body>
    </html>
  )
}
