import { Button } from '@/components/ui/button';
import * as Constants from '@/constants';
import axios from 'axios';
import { ChangeEvent, FormEvent, useRef, useState } from 'react';

type FormData = {
  email: string;
  id: string;
  data: any;
  ts: string;
};

const ContactForm = () => {
  const [isSubmitted, setIsSubmitted] = useState(false);
  const timeoutRef = useRef(null);

  const [formData, setFormData] = useState({
    email: '',
    message: ''
  });

  const handleInputChange = (
    e: ChangeEvent<HTMLInputElement | HTMLTextAreaElement>
  ) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));
  };

  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    const reqBody: FormData = {
      email: formData.email,
      data: formData.message,
      id: Constants.APP_TITLE,
      ts: (new Date()).toISOString(),
    }

    try {
      const response = await axios.put("https://" + Constants.FORMS_HOST + "/v1/entries/", reqBody)

      if (response.status === 200) {
        setIsSubmitted(true);
      }
    } catch (error) {
      console.error(error);
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <div className="flex items-center justify-between">
        <h3 className="text-2xl font-heading">
          Entre em Contato
        </h3>
        {isSubmitted && (
          <span className="border-border text-text dark:border-darkBorder rounded-base border-2 bg-main px-2 py-0.5 text-sm">
            Enviado!
          </span>
        )}
      </div>
      <p className="mb-3 mt-1">
        Entraremos em contato para uma demonstração rápida
      </p>
      <div className="space-y-4">
        <label className="text-2xl font-heading">
          email
          <input
            type="email"
            name="email"
            value={formData.email}
            // disabled={isSubmitted}
            onChange={handleInputChange}
            className="border-border text-text dark:border-darkBorder rounded-base w-full border-2 px-2 py-0.5 text-sm font-light"
          />
        </label>
        <label className="text-2xl font-heading">
          mensagem
          <textarea
            name="message"
            value={formData.message}
            // disabled={isSubmitted}
            onChange={handleInputChange}
            className="h-40 border-border text-text dark:border-darkBorder rounded-base w-full border-2 px-2 py-0.5 text-sm font-light"
          />
        </label>
      </div>
      <Button
        type="submit"
        size="lg"
        className="mt-6 w-full"
      >
        Enviar
      </Button>
    </form>
  );
};

export default ContactForm;