'use strict';

module.exports = function (plop) {
  plop.setGenerator('module', {
    description: 'Scaffold a new Bridge API module (model + controller + seeder + routes)',

    prompts: [
      {
        type: 'input',
        name: 'name',
        message: 'Module name (PascalCase, e.g. Driver):',
        validate: (v) =>
          /^[A-Z][a-zA-Z0-9]+$/.test(v) || 'Must be PascalCase and a single word (e.g. Driver)',
      },
      {
        type: 'input',
        name: 'plural',
        message: 'Route path — plural lowercase (e.g. drivers):',
        default: (a) => a.name.toLowerCase() + 's',
      },
    ],

    actions: [
      // ── New files ──────────────────────────────────────────────────────────

      {
        type: 'add',
        path: 'internal/models/{{name}}.go',
        templateFile: 'plop-templates/model.go.hbs',
      },
      {
        type: 'add',
        path: 'internal/controllers/{{camelCase name}}/{{camelCase name}}_controller.go',
        templateFile: 'plop-templates/controller.go.hbs',
      },
      {
        type: 'add',
        path: 'internal/database/seeder/{{camelCase name}}_form_fields.go',
        templateFile: 'plop-templates/seeder.go.hbs',
      },

      // ── Inject into routes.go ──────────────────────────────────────────────

      // 1. Add import
      {
        type: 'modify',
        path: 'internal/routes/routes.go',
        pattern: /(\/\/ plop:imports)/,
        template:
          '\t{{camelCase name}} "github.com/AhmadAboElzahab/bridge/internal/controllers/{{camelCase name}}"\n\t$1',
      },
      // 2. Add controller instantiation
      {
        type: 'modify',
        path: 'internal/routes/routes.go',
        pattern: /(\/\/ plop:controllers)/,
        template: '\t{{camelCase name}}Ctrl := {{camelCase name}}.New{{name}}Controller()\n\t$1',
      },
      // 3. Add route group
      {
        type: 'modify',
        path: 'internal/routes/routes.go',
        pattern: /(\/\/ plop:routes)/,
        template: [
          '\t{{plural}}Routes := protected.Group("/{{plural}}")',
          '\t{',
          '\t\t{{plural}}Routes.POST("/index", {{camelCase name}}Ctrl.Index)',
          '\t\t{{plural}}Routes.POST("/", {{camelCase name}}Ctrl.Store)',
          '\t\t{{plural}}Routes.GET("/:id", {{camelCase name}}Ctrl.Show)',
          '\t\t{{plural}}Routes.PUT("/:id", {{camelCase name}}Ctrl.Update)',
          '\t\t{{plural}}Routes.DELETE("/:id", {{camelCase name}}Ctrl.Delete)',
          '\t}',
          '\t$1',
        ].join('\n'),
      },

      // ── Inject into migrations/main.go ────────────────────────────────────

      {
        type: 'modify',
        path: 'cmd/migrations/main.go',
        pattern: /(\/\/ plop:models)/,
        template: 'initializers.DB.AutoMigrate(&models.{{name}}{})\n\t$1',
      },

      // ── Inject into seeder/runner.go ──────────────────────────────────────

      // Add individual case
      {
        type: 'modify',
        path: 'internal/database/seeder/runner.go',
        pattern: /(\/\/ plop:cases)/,
        template: 'case "{{plural}}":\n\t\tSeed{{name}}FormFields()\n\t$1',
      },
      // Add to "all" case
      {
        type: 'modify',
        path: 'internal/database/seeder/runner.go',
        pattern: /(\/\/ plop:all)/,
        template: 'Seed{{name}}FormFields()\n\t\t$1',
      },
    ],
  });
};
