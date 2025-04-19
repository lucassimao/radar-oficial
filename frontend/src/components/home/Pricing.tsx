import React from "react";

interface FeatureProps {
  available: boolean;
  text: string;
}

const Feature = ({ available, text }: FeatureProps) => (
  <div className="flex items-center py-2">
    <div className="mr-3 text-xl">
      {available ? (
        <span className="text-green-500">✓</span>
      ) : (
        <span className="text-gray-400">❌</span>
      )}
    </div>
    <span className={available ? "text-gray-800" : "text-gray-500"}>
      {text}
    </span>
  </div>
);

interface PricingCardProps {
  plan: string;
  price: string;
  limit: string;
  features: {
    text: string;
    available: boolean;
  }[];
  highlighted?: boolean;
}

const PricingCard = ({
  plan,
  price,
  limit,
  features,
  highlighted = false,
}: PricingCardProps) => (
  <div
    className={`rounded-2xl shadow-lg ${highlighted ? "border-2 border-blue-500 transform scale-105 z-10" : ""} overflow-hidden bg-white transition-all duration-300 hover:shadow-xl hover:scale-[1.02]`}
  >
    {highlighted && (
      <div className="bg-blue-600 text-white text-center py-2 font-medium">
        Recomendado
      </div>
    )}
    <div className="p-8">
      <h3 className="text-2xl font-bold mb-2">{plan}</h3>
      <div className="mb-4">
        <span className="text-4xl font-bold">{price}</span>
        {price !== "Gratuito" && <span className="text-gray-600">/mês</span>}
      </div>
      <p className="text-gray-600 mb-6 pb-4 border-b border-gray-100">
        {limit}
      </p>

      <div className="space-y-1 mb-8">
        {features.map((feature, index) => (
          <Feature
            key={index}
            available={feature.available}
            text={feature.text}
          />
        ))}
      </div>

      <button className="w-full mt-4 py-3 px-4 rounded-lg bg-gray-200 text-gray-500 font-medium cursor-not-allowed opacity-70 transition-all duration-200">
        Disponível em breve
      </button>
    </div>
  </div>
);

const Pricing = () => {
  const plans = [
    {
      plan: "Gratuito",
      price: "Gratuito",
      limit: "5 buscas/mês",
      highlighted: false,
      features: [
        { text: "Busca inteligente", available: true },
        { text: "Interface conversacional (Web/WhatsApp)", available: true },
        {
          text: "Notificações por e-mail com palavras-chave",
          available: false,
        },
        { text: "Acesso ao histórico de buscas", available: false },
      ],
    },
    {
      plan: "Básico",
      price: "R$ 19,90",
      limit: "5 buscas/dia",
      highlighted: true,
      features: [
        { text: "Busca inteligente", available: true },
        { text: "Interface conversacional (Web/WhatsApp)", available: true },
        {
          text: "Notificações por e-mail com palavras-chave",
          available: false,
        },
        { text: "Acesso ao histórico de buscas", available: false },
      ],
    },
    {
      plan: "Profissional",
      price: "R$ 49,90",
      limit: "30 buscas/dia",
      highlighted: false,
      features: [
        { text: "Busca inteligente", available: true },
        { text: "Interface conversacional (Web/WhatsApp)", available: true },
        { text: "Notificações por e-mail com palavras-chave", available: true },
        { text: "Acesso ao histórico de buscas", available: true },
      ],
    },
  ];

  return (
    <section id="pricing" className="py-20 bg-gray-50">
      <div className="max-w-screen-lg mx-auto px-4">
        <div className="text-center mb-16">
          <h2 className="text-3xl md:text-4xl font-bold mb-4">
            Escolha o plano ideal para você
          </h2>
          <p className="text-xl text-gray-600 max-w-3xl mx-auto">
            Opções flexíveis para diferentes necessidades de monitoramento de
            publicações oficiais.
          </p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-4 md:gap-8 max-w-6xl mx-auto">
          {plans.map((plan, index) => (
            <PricingCard
              key={index}
              plan={plan.plan}
              price={plan.price}
              limit={plan.limit}
              features={plan.features}
              highlighted={plan.highlighted}
            />
          ))}
        </div>

        <div className="mt-12 text-center text-sm text-gray-500">
          <p>
            Preços de pré-lançamento. Sujeitos a alterações até a data de
            lançamento oficial.
          </p>
        </div>
      </div>
    </section>
  );
};

export default Pricing;
