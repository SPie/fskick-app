<?php

namespace App\Models;

use Doctrine\ORM\Mapping;

/**
 * Class AbstractDoctrineModel
 *
 * @package App\Models
 */
abstract class AbstractDoctrineModel implements Model
{
    /**
     * @Mapping\Id
     * @Mapping\GeneratedValue
     * @Mapping\Column(type="integer")
     *
     * @var int|null
     */
    protected ?int $id;

    /**
     * @param int|null $id
     *
     * @return $this|Model
     */
    public function setId(?int $id): Model
    {
        $this->id = $id;

        return $this;
    }

    /**
     * @return int|null
     */
    public function getId(): ?int
    {
        return $this->id;
    }
}
