
import type { NextPage } from "next";
import Head from "next/head";
import { BoltIcon, MagnifyingGlassIcon, ChatBubbleLeftRightIcon } from '@heroicons/react/24/outline';

const Home: NextPage = () => {



  return (
    <div className="min-h-screen bg-gray-50">
      <Head>
        <title>Radar Oficial - Diários Oficiais Inteligentes</title>
        <meta name="description" content="Acompanhe os Diários Oficiais com inteligência" />
        <link rel="icon" href="/favicon.ico" />
      </Head>

      {/* Hero Section */}
      <section className="relative bg-white overflow-hidden">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-20">
          <div className="text-center">
            <span className="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-blue-50 text-accent mb-6">
              🚀 Em breve
            </span>
            <h1 className="text-4xl sm:text-5xl font-bold text-gray-900 mb-6">
              Acompanhe os Diários Oficiais com inteligência
            </h1>
            <p className="text-xl text-gray-600 mb-8 max-w-3xl mx-auto">
              Radar Oficial reúne os Diários Oficiais do Piauí em um só lugar — com busca inteligente e acesso instantâneo via WhatsApp.
            </p>
            <form  className="max-w-md mx-auto mb-4">
              <div className="flex gap-2">
                <input
                  type="email"
                  placeholder="Seu email"
                  className="flex-1 px-4 py-3 rounded-lg border border-gray-300 focus:outline-none focus:ring-2 focus:ring-accent"
                  required
                />
                <button type="submit" className="px-6 py-3 bg-accent text-white rounded-lg font-medium hover:bg-blue-700 transition">
                  Avise-me no lançamento
                </button>
              </div>
            </form>
            <p className="text-gray-500">Lançamento previsto para maio de 2025</p>
          </div>
        </div>
      </section>

      {/* Features Section */}
      <section className="py-20 bg-white">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="grid grid-cols-1 md:grid-cols-3 gap-12">
            <div className="text-center">
              <div className="inline-flex items-center justify-center w-16 h-16 rounded-lg bg-blue-50 mb-6">
                <MagnifyingGlassIcon className="h-8 w-8 text-accent" />
              </div>
              <h3 className="text-xl font-semibold mb-3">Consulta Inteligente</h3>
              <p className="text-gray-600">Pesquise termos e entidades públicas com precisão.</p>
            </div>
            <div className="text-center">
              <div className="inline-flex items-center justify-center w-16 h-16 rounded-lg bg-blue-50 mb-6">
                <BoltIcon className="h-8 w-8 text-accent" />
              </div>
              <h3 className="text-xl font-semibold mb-3">Notificações em Tempo Real</h3>
              <p className="text-gray-600">Receba alertas sempre que novos diários forem publicados.</p>
            </div>
            <div className="text-center">
              <div className="inline-flex items-center justify-center w-16 h-16 rounded-lg bg-blue-50 mb-6">
                <ChatBubbleLeftRightIcon className="h-8 w-8 text-accent" />
              </div>
              <h3 className="text-xl font-semibold mb-3">Assistente no WhatsApp</h3>
              <p className="text-gray-600">Consulte diretamente os diários com linguagem natural.</p>
            </div>
          </div>
        </div>
      </section>

      {/* Pricing Section */}
      <section className="py-20 bg-gray-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <h2 className="text-3xl font-bold text-center mb-12">Planos</h2>
          <p className="text-center text-gray-600 mb-8">Planos disponíveis após o lançamento</p>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
            {[
              { name: "Gratuito", price: "R$ 0", limit: "até 5 buscas/mês" },
              { name: "Básico", price: "R$ 19,90/mês", limit: "5 buscas/dia", featured: true },
              { name: "Pro", price: "R$ 49,90/mês", limit: "30 buscas/dia" }
            ].map((plan) => (
              <div 
                key={plan.name} 
                className={`bg-white rounded-xl p-8 ${plan.featured ? 'ring-2 ring-accent shadow-lg' : 'border border-gray-200'}`}
              >
                <h3 className="text-xl font-semibold mb-2">{plan.name}</h3>
                <p className="text-3xl font-bold mb-4">{plan.price}</p>
                <p className="text-gray-600 mb-6">{plan.limit}</p>
                <button className={`w-full py-3 px-6 rounded-lg font-medium transition ${plan.featured ? 'bg-accent text-white hover:bg-blue-700' : 'bg-gray-100 text-gray-800 hover:bg-gray-200'}`}>
                  Avise-me no lançamento
                </button>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Footer */}
      <footer className="bg-white border-t">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
          <div className="flex flex-col md:flex-row justify-between items-center gap-4">
            <div className="flex space-x-6">
              <a href="#" className="text-gray-600 hover:text-gray-900">Termos de Uso</a>
              <a href="#" className="text-gray-600 hover:text-gray-900">Privacidade</a>
              <a href="#" className="text-gray-600 hover:text-gray-900">Contato</a>
            </div>
            <p className="text-gray-600">© 2024 Radar Oficial</p>
          </div>
        </div>
      </footer>
    </div>
  );
};

export default Home;
