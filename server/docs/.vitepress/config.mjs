import { defineConfig } from 'vitepress'

export default defineConfig({
  title: 'Bridge API',
  description: 'UAE domestic service marketplace — API documentation',

  // Read from /docs, ignore Swagger-generated files
  srcExclude: ['**/swagger.json', '**/swagger.yaml', '**/docs.go'],

  // localhost links are only reachable when the API is running locally — ignore during build
  ignoreDeadLinks: [/localhost/],

  themeConfig: {
    logo: '🌉',
    siteTitle: 'Bridge API',

    nav: [
      { text: 'Guide', link: '/getting-started' },
      { text: 'API Reference', link: '/api-reference' },
      {
        text: 'Swagger UI',
        link: 'http://localhost:8080/swagger/index.html',
        target: '_blank'
      }
    ],

    sidebar: [
      {
        text: 'Introduction',
        items: [
          { text: 'Overview', link: '/' },
          { text: 'Getting Started', link: '/getting-started' }
        ]
      },
      {
        text: 'Configuration',
        items: [
          { text: 'Environment Variables', link: '/environment-variables' },
          { text: 'Storage (Local / S3 / R2)', link: '/storage' }
        ]
      },
      {
        text: 'Reference',
        items: [
          { text: 'Architecture', link: '/architecture' },
          { text: 'API Reference', link: '/api-reference' },
          { text: 'Tab & Filter System', link: '/tab-filter-system' }
        ]
      },
      {
        text: 'Operations',
        items: [
          { text: 'Deployment', link: '/deployment' }
        ]
      }
    ],

    search: {
      provider: 'local'
    },

    socialLinks: [
      { icon: 'github', link: 'https://github.com/AhmadAboElzahab/bridge' }
    ],

    footer: {
      message: 'Bridge API — Internal Documentation',
      copyright: 'UAE Domestic Service Marketplace'
    },

    editLink: {
      pattern: 'https://github.com/AhmadAboElzahab/bridge/edit/main/docs/:path',
      text: 'Edit this page on GitHub'
    }
  }
})
