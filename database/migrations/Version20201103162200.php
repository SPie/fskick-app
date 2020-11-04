<?php

namespace Database\Migrations;

use Doctrine\DBAL\Schema\Schema;
use Doctrine\Migrations\AbstractMigration;
use LaravelDoctrine\Migrations\Schema\Builder;
use LaravelDoctrine\Migrations\Schema\Table;

/**
 * Class Version20201103162200
 *
 * @package Database\Migrations
 */
final class Version20201103162200 extends AbstractMigration
{
    /**
     * @param Schema $schema
     *
     * @return void
     */
    public function up(Schema $schema): void
    {
        $this
            ->createPlayersTable($schema)
            ->createSeasonsTable($schema);
    }

    /**
     * @param Schema $schema
     *
     * @return $this
     */
    private function createPlayersTable(Schema $schema): self
    {
        (new Builder($schema))->create('players', function (Table $table) {
            $table->increments('id');
            $table->string('name');
            $table->unique('name');
            $table->timestamps();
        });

        return $this;
    }

    /**
     * @param Schema $schema
     *
     * @return $this
     */
    private function createSeasonsTable(Schema $schema): self
    {
        (new Builder($schema))->create('seasons', function (Table $table) {
            $table->increments('id');
            $table->string('name');
            $table->unique('name');
            $table->boolean('active')->setDefault(false);
            $table->timestamps();;
        });

        return $this;
    }

    /**
     * @param Schema $schema
     *
     * @return void
     */
    public function down(Schema $schema): void
    {
        $this
            ->dropSeasonsTable($schema)
            ->dropPlayersTable($schema);
    }

    /**
     * @param Schema $schema
     *
     * @return $this
     */
    private function dropPlayersTable(Schema $schema): self
    {
        (new Builder($schema))->dropIfExists('players');

        return $this;
    }

    /**
     * @param Schema $schema
     *
     * @return $this
     */
    private function dropSeasonsTable(Schema $schema): self
    {
        (new Builder($schema))->dropIfExists('seasons');

        return $this;
    }
}
