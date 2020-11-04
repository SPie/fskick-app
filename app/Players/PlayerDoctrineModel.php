<?php

namespace App\Players;

use App\Models\AbstractDoctrineModel;
use App\Models\Timestamps;
use Doctrine\ORM\Mapping as ORM;

/**
 * Class PlayerDoctrineModel
 *
 * @ORM\Table(name="players")
 * @ORM\Entity(repositoryClass="App\Players\PlayerDoctrineRepository")
 *
 * @package App\Players
 */
final class PlayerDoctrineModel extends AbstractDoctrineModel implements PlayerModel
{
    use Timestamps;

    /**
     * @ORM\Column(name="name", type="string", length=255, nullable=false)
     *
     * @var string
     */
    private string $name;

    /**
     * PlayerDoctrineModel constructor.
     *
     * @param string $name
     */
    public function __construct(string $name)
    {
        $this->name = $name;
    }

    /**
     * @return string
     */
    public function getName(): string
    {
        return $this->name;
    }
}
