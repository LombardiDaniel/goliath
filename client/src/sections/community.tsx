'use client'

import ContactForm from '@/components/contact-form';


export default function Community() {
  return (
    <section id="community" className="border-b-border dark:border-b-darkBorder dark:bg-secondaryBlack inset-0 flex w-full flex-col items-center justify-center border-b-2 bg-white bg-[linear-gradient(to_right,#80808033_1px,transparent_1px),linear-gradient(to_bottom,#80808033_1px,transparent_1px)] bg-[size:70px_70px] font-base">
      <div className="mx-auto w-container max-w-full px-5 py-20 lg:py-[100px]">
        <h2 className="mb-14 text-center text-2xl font-heading md:text-3xl lg:mb-20 lg:text-4xl">
          Agende uma demonstração
        </h2>

        <div className="items-center grid grid-cols-1 gap-8 w900:mx-auto w900:w-2/3 w900:grid-cols-1 w500:w-full">
          <div className="border-border dark:border-darkBorder dark:bg-secondaryBlack flex flex-col justify-between rounded-base border-2 bg-white p-5">
            <ContactForm/>
          </div>
        </div>
      </div>
    </section>
  )
}
