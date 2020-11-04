<?php

namespace App\Seasons;

use App\Models\AbstractDoctrineModel;
use App\Models\Timestamps;
use Doctrine\ORM\Mapping as ORM;

/**
 * Class SeasonDoctrineModel
 *
 * @ORM\Table(name="seasons")
 * @ORM\Entity(repositoryClass="App\Seasons\SeasonDoctrineRepository")
 *
 * @package App\Seasons
 */
final class SeasonDoctrineModel extends AbstractDoctrineModel implements SeasonModel
{
    use Timestamps;

    /**
     * @ORM\Column(name="name", type="string", length=255, nullable=false)
     *
     * @var string
     */
    private string $name;

    /**
     * @ORM\Column(name="active", type="boolean")
     *
     * @var bool
     */
    private bool $active;

    /**
     * SeasonDoctrineModel constructor.
     *
     * @param string $name
     * @param bool   $active
     */
    public function __construct(string $name, bool $active = false)
    {
        $this->name = $name;
        $this->active = $active;
    }

    /**
     * @return string
     */
    public function getName(): string
    {
        return $this->name;
    }

    /**
     * @param bool $active
     *
     * @return SeasonModel
     */
    public function setActive(bool $active): SeasonModel
    {
        $this->active = $active;

        return $this;
    }

    /**
     * @return bool
     */
    public function isActive(): bool
    {
        return $this->active;
    }
}
