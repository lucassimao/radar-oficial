import React from 'react';

const FAQItem = ({ 
  question, 
  answer 
}: { 
  question: string, 
  answer: string 
}) => (
  <div className="border-b border-gray-200 pb-6 mb-6 last:border-0 last:mb-0 last:pb-0">
    <h3 className="text-xl font-bold mb-3 text-gray-900">{question}</h3>
    <p className="text-gray-600">{answer}</p>
  </div>
);

const FAQ = () => {
  const faqs = [
    {
      question: "Quais Diários Oficiais estão disponíveis no Radar Oficial?",
      answer: "O Radar Oficial inclui o Diário Oficial do Estado do Piauí e os Diários Oficiais de todos os municípios piauienses, além de publicações de autarquias, fundações e outros órgãos públicos do estado."
    },
    {
      question: "Com que frequência as publicações são atualizadas?",
      answer: "As publicações são atualizadas em tempo real, assim que são disponibilizadas pelos órgãos oficiais. Você receberá notificações imediatas sobre novas publicações que correspondam aos seus alertas."
    },
    {
      question: "Posso acessar o Radar Oficial pelo celular?",
      answer: "Sim, o Radar Oficial é responsivo e funciona em qualquer dispositivo. Além disso, oferecemos acesso via WhatsApp, permitindo que você consulte publicações e receba alertas diretamente no seu celular."
    },
    {
      question: "Como funcionam os alertas personalizados?",
      answer: "Você pode configurar alertas com base em palavras-chave, nomes, CPF/CNPJ, órgãos específicos ou qualquer outro termo relevante. Sempre que uma publicação corresponder aos seus critérios, você receberá uma notificação por e-mail ou WhatsApp."
    },
    {
      question: "O Radar Oficial pode ser usado por órgãos públicos?",
      answer: "Sim, oferecemos planos específicos para órgãos públicos, com recursos adicionais de monitoramento e relatórios customizados. Entre em contato para saber mais sobre nosso plano Empresarial."
    },
    {
      question: "Quando o Radar Oficial estará disponível?",
      answer: "O lançamento oficial está previsto para maio de 2025. Estamos trabalhando arduamente para entregar uma plataforma robusta e confiável para o monitoramento de publicações oficiais no Piauí."
    }
  ];

  return (
    <section id="faq" className="py-20 bg-gray-50">
      <div className="container mx-auto px-4">
        <div className="text-center mb-16">
          <h2 className="text-3xl md:text-4xl font-bold mb-4">Perguntas Frequentes</h2>
          <p className="text-xl text-gray-600 max-w-3xl mx-auto">
            Tire suas dúvidas sobre o Radar Oficial e nossos serviços.
          </p>
        </div>
        
        <div className="max-w-4xl mx-auto bg-white rounded-xl shadow-lg p-8">
          {faqs.map((faq, index) => (
            <FAQItem 
              key={index}
              question={faq.question}
              answer={faq.answer}
            />
          ))}
        </div>
      </div>
    </section>
  );
};

export default FAQ;
