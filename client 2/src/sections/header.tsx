import { Button } from '@/components/ui/button';
import * as Constants from "@/constants";
import { Calendar, HandCoins } from 'lucide-react';
import { Fragment } from 'react';

export default function Header() {
  return (
    <header id="header" className="dark:bg-secondaryBlack inset-0 flex min-h-[80dvh] w-full flex-col items-center justify-center bg-white bg-[linear-gradient(to_right,#80808033_1px,transparent_1px),linear-gradient(to_bottom,#80808033_1px,transparent_1px)] bg-[size:70px_70px]">
      <div className="mx-auto w-container max-w-full px-5 py-[110px] text-center lg:py-[150px]">
        <h1 className="text-3xl font-heading md:text-4xl lg:text-5xl">
          {Constants.APP_TITLE}
        </h1>
        <p className="my-12 mt-8 text-lg font-normal leading-relaxed md:text-xl lg:text-2xl lg:leading-relaxed">
          {Constants.DESCRIPTION.map((item, index) => (
            <Fragment key={index}>
              {item}
              {index < Constants.DESCRIPTION.length - 1 && <br />}
            </Fragment>
          ))}
        </p>
        <div className="flex flex-col md:flex-row justify-center space-y-4 md:space-y-0 md:space-x-4 w-fit mx-auto">
          <a href="#pricing">
            <Button
              size="lg"
              className="w-full h-12 text-base font-heading md:text-lg lg:h-14 lg:text-xl"
            >
              <HandCoins className="mr-2 h-5 w-5" strokeWidth={3} />
              Confira nossos planos
            </Button>
          </a>

          <a href="#community">
            <Button
              size="lg"
              className="w-full h-12 text-base font-heading md:text-lg lg:h-14 lg:text-xl"
            >
              <Calendar className="mr-2 h-5 w-5" strokeWidth={3} />
              Agende uma demonstração
            </Button>
          </a>
        </div>
      </div>
    </header>
  )
}
